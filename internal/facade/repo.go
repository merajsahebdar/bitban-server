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
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/server"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"regeet.io/api/internal/conf"
)

// Repo
type Repo struct {
	ctx        context.Context
	fs         billy.Filesystem
	storage    storage.Storer
	loader     server.Loader
	repository *git.Repository
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

// advertiseRefs
func (f *Repo) advertiseRefs(w io.Writer, sess transport.Session) error {
	ar, err := sess.AdvertisedReferencesContext(f.ctx)
	if err != nil {
		return err
	}

	return ar.Encode(w)
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

// AdvertiseRefs
func (f *Repo) AdvertiseRefs(service string, w io.Writer) error {
	var err error

	var sess transport.Session

	switch service {
	case transport.ReceivePackServiceName:
		sess, err = f.initReceivePackSession()
	case transport.UploadPackServiceName:
		sess, err = f.initUploadPackSession()
	}

	if err != nil {
		return err
	}

	if ar, err := sess.AdvertisedReferencesContext(f.ctx); err != nil {
		return err
	} else {
		enc := pktline.NewEncoder(w)
		enc.Encodef("# service=%s\n", service)
		enc.Flush()

		return ar.Encode(w)
	}
}

// ReceivePack
func (f *Repo) ReceivePack(r io.Reader, w io.Writer, adv bool) error {
	sess, err := f.initReceivePackSession()
	if err != nil {
		return err
	}

	req := packp.NewReferenceUpdateRequest()

	if adv {
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
}

// UploadPack
func (f *Repo) UploadPack(r io.Reader, w io.Writer, adv bool) error {
	sess, err := f.initUploadPackSession()
	if err != nil {
		return err
	}

	req := packp.NewUploadPackRequest()

	if adv {
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

// newStorage
func newStorage(name string) (billy.Filesystem, storage.Storer) {
	fs := osfs.New(
		os.Getenv("HOME") + conf.Cog.Storage.Dir + "/" + name,
	)

	return fs, filesystem.NewStorage(
		fs,
		cache.NewObjectLRUDefault(),
	)
}

// GetRepoByName
func GetRepoByName(ctx context.Context, name string) (*Repo, error) {
	fs, storage := newStorage(name)
	if repository, err := git.Open(storage, nil); err != nil {
		return nil, err
	} else {
		return &Repo{
			ctx:        ctx,
			fs:         fs,
			storage:    storage,
			loader:     &serverLoader{storage: storage},
			repository: repository,
		}, nil
	}
}
