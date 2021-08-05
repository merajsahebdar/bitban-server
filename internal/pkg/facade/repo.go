/*
 * Copyright 2021 Meraj Sahebdar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package facade

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/uptrace/bun"
	"github.com/volatiletech/null/v8"
	gossh "golang.org/x/crypto/ssh"
	"regeet.io/api/internal/cfg"
	"regeet.io/api/internal/pkg/exec"
	"regeet.io/api/internal/pkg/orm"
	"regeet.io/api/internal/pkg/orm/entity"
)

// Services
const (
	GitReceivePack = "git-receive-pack"
	GitUploadPack  = "git-upload-pack"
)

// Capabilities
const (
	NoThin capability.Capability = "no-thin"
)

// ServerPackConfig
type ServerPackConfig struct {
	R       io.Reader
	W       io.Writer
	Service string
	IsSsh   bool
}

// repoBackend
type repoGoBackend struct {
	fs         billy.Filesystem
	storage    storage.Storer
	loader     server.Loader
	repository *git.Repository
}

// Repo
type Repo struct {
	*repoGoBackend
	ctx           context.Context
	domainAddress string
	repoAddress   string
	path          string
}

type serverLoader struct {
	storage storage.Storer
}

// Load
func (l *serverLoader) Load(ep *transport.Endpoint) (storer.Storer, error) {
	return l.storage, nil
}

// getTransportServer
func (f *Repo) getTransportServer() transport.Transport {
	return server.NewServer(f.loader)
}

// initReceivePackSession
func (f *Repo) initReceivePackSession() (transport.ReceivePackSession, error) {
	sess, err := f.getTransportServer().NewReceivePackSession(&transport.Endpoint{}, nil)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// initUploadPackSession
func (f *Repo) initUploadPackSession() (transport.UploadPackSession, error) {
	sess, err := f.getTransportServer().NewUploadPackSession(&transport.Endpoint{}, nil)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// setAdvertisingCapabilites
func setAdvertisingCapabilites(ar *packp.AdvRefs) {
	ar.Capabilities.Add(NoThin)
}

// advertiseRefs
func (f *Repo) advertiseRefs(w io.Writer, sess transport.Session) error {
	ar, err := sess.AdvertisedReferencesContext(f.ctx)
	if err != nil {
		return err
	}

	setAdvertisingCapabilites(ar)

	return ar.Encode(w)
}

// AdvertiseRefs
func (f *Repo) AdvertiseRefs(w io.Writer, service string) error {
	if cfg.IsGoBackend() {
		var err error

		var sess transport.Session

		switch service {
		case GitReceivePack:
			sess, err = f.initReceivePackSession()
		case GitUploadPack:
			sess, err = f.initUploadPackSession()
		}

		if err != nil {
			return err
		}

		if ar, err := sess.AdvertisedReferencesContext(f.ctx); err != nil {
			return err
		} else {
			setAdvertisingCapabilites(ar)

			enc := pktline.NewEncoder(w)
			enc.Encodef("# service=%s\n", service)
			enc.Flush()

			return ar.Encode(w)
		}
	} else {
		if cmd, _, stdout, stderr, err := exec.Create(
			"git",
			strings.TrimPrefix(service, "git-"),
			"--stateless-rpc",
			"--advertise-refs",
			".",
		); err != nil {
			return err
		} else {
			cmd.Dir = f.path

			if err := cmd.Start(); err != nil {
				return err
			}

			enc := pktline.NewEncoder(w)
			enc.Encodef("# service=%s\n", service)

			io.Copy(w, stdout)
			io.Copy(w, stderr)

			return cmd.Wait()
		}
	}
}

// ServePack
func (f *Repo) ServePack(serveConfig *ServerPackConfig) error {
	r := serveConfig.R
	w := serveConfig.W

	if cfg.IsGoBackend() {
		//
		// Serve by go-git package.

		switch serveConfig.Service {
		case GitReceivePack:
			sess, err := f.initReceivePackSession()
			if err != nil {
				return err
			}

			req := packp.NewReferenceUpdateRequest()

			if serveConfig.IsSsh {
				if err = f.advertiseRefs(w, sess); err != nil {
					return err
				}
			}

			if err := req.Decode(r); err != nil {
				return err
			}

			if status, err := sess.ReceivePack(f.ctx, req); status != nil {
				return status.Encode(w)
			} else {
				return err
			}
		case GitUploadPack:
			sess, err := f.initUploadPackSession()
			if err != nil {
				return err
			}

			req := packp.NewUploadPackRequest()

			if serveConfig.IsSsh {
				if err = f.advertiseRefs(w, sess); err != nil {
					return err
				}
			}

			if err := req.Decode(r); err != nil {
				return err
			}

			if status, err := sess.UploadPack(f.ctx, req); status != nil {
				return status.Encode(w)
			} else {
				return err
			}
		}

		panic("not a valid git service")
	} else {
		//
		// Serve by git binary.

		args := []string{
			strings.TrimPrefix(serveConfig.Service, "git-"),
		}

		if !serveConfig.IsSsh {
			args = append(args, "--stateless-rpc")
		}

		args = append(args, ".")

		if cmd, stdin, stdout, stderr, err := exec.Create(
			"git",
			args...,
		); err != nil {
			return err
		} else {
			cmd.Dir = f.path

			if err := cmd.Start(); err != nil {
				return err
			}

			go io.Copy(stdin, r)
			io.Copy(w, stdout)
			if ch, ok := w.(gossh.Channel); ok {
				io.Copy(ch.Stderr(), stderr)
			}

			return cmd.Wait()
		}
	}
}

// CreateRepoByAddress
func CreateRepoByAddress(ctx context.Context, domainAddress string, repoAddress string) (repo *Repo, err error) {
	domain := new(entity.Domain)
	if err := orm.
		GetBunInstance().
		NewSelect().
		Model(domain).
		Where("? = ?", bun.Ident("domain.address"), domainAddress).
		Scan(ctx); err != nil {
		return nil, err
	}

	var tx bun.Tx
	if tx, err = orm.GetBunInstance().BeginTx(ctx, nil); err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && repo == nil {
			tx.Rollback()
		}
	}()

	repository := &entity.Repository{
		Address:  repoAddress,
		DomainID: null.Int64From(domain.ID),
	}
	if _, err := tx.
		NewInsert().
		Model(repository).
		Column("address", "domain_id").
		Exec(ctx); err != nil {
		return nil, err
	}

	path, err := getPath(domainAddress, repoAddress)
	if err != nil {
		return nil, err
	}

	var backend *repoGoBackend
	if cfg.IsGoBackend() {
		fs, storage, err := newStorage(path)
		if err != nil {
			return nil, err
		}

		if repository, err := git.Init(storage, nil); err != nil {
			return nil, err
		} else {
			head := plumbing.NewSymbolicReference(
				plumbing.HEAD,
				plumbing.ReferenceName(fmt.Sprintf(
					"refs/heads/%s",
					cfg.Cog.Git.Configs.Init.DefaultBranch,
				)),
			)
			if err := storage.SetReference(head); err != nil {
				return nil, err
			}

			backend = &repoGoBackend{
				fs:         fs,
				storage:    storage,
				loader:     &serverLoader{storage: storage},
				repository: repository,
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	repo = &Repo{
		repoGoBackend: backend,
		ctx:           ctx,
		domainAddress: domainAddress,
		repoAddress:   repoAddress,
		path:          path,
	}

	return repo, nil
}

// GetRepoByAddress
func GetRepoByAddress(ctx context.Context, domainAddress string, repoAddress string) (*Repo, error) {
	repository := new(entity.Repository)
	if err := orm.
		GetBunInstance().
		NewSelect().
		Model(repository).
		Relation("Domain", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("? = ?", bun.Ident("domain.address"), domainAddress)
		}).
		Where("? = ?", bun.Ident("repository.address"), repoAddress).
		Scan(ctx); err != nil {
		return nil, err
	}

	path, err := getPath(domainAddress, repoAddress)
	if err != nil {
		return nil, err
	}

	var backend *repoGoBackend
	if cfg.IsGoBackend() {
		fs, storage, err := newStorage(path)
		if err != nil {
			return nil, err
		}

		if repository, err := git.Open(storage, nil); err != nil {
			return nil, err
		} else {
			backend = &repoGoBackend{
				fs:         fs,
				storage:    storage,
				loader:     &serverLoader{storage: storage},
				repository: repository,
			}
		}
	}

	return &Repo{
		repoGoBackend: backend,
		ctx:           ctx,
		domainAddress: domainAddress,
		repoAddress:   repoAddress,
		path:          path,
	}, nil
}

// getPath
func getPath(domainAddress string, repoAddress string) (path string, err error) {
	if cfg.Cog.Git.Storage == cfg.GitStorageFs {
		if path, err = cfg.GetVarPath("/repos", domainAddress, repoAddress); err != nil {
			return "", git.ErrRepositoryNotExists
		}
	} else {
		path = fmt.Sprintf("/repos/%s/%s", domainAddress, repoAddress)
	}

	return path, nil
}

// newStorage
func newStorage(path string) (billy.Filesystem, storage.Storer, error) {
	if fs, err := getFs(path); err != nil {
		return nil, nil, err
	} else {
		return fs, filesystem.NewStorage(
			fs,
			cache.NewObjectLRUDefault(),
		), nil
	}
}

// fsMutex
var fsMutex = &sync.Mutex{}

// fs
var fs billy.Filesystem

// getFs
func getFs(path string) (billy.Filesystem, error) {
	if fs == nil {
		fsMutex.Lock()
		defer fsMutex.Unlock()

		if fs == nil {
			switch cfg.Cog.Git.Storage {
			case cfg.GitStorageFs:
				if p, err := cfg.GetVarPath("/repos"); err != nil {
					return nil, err
				} else {
					fs = osfs.New(p)
				}
			case cfg.GitStorageMem:
				fs = memfs.New()
			}
		}
	}

	return fs.Chroot(path)
}
