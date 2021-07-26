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

func (l *serverLoader) Load(ep *transport.Endpoint) (storer.Storer, error) {
	return l.storage, nil
}

// AdvertiseRefs
func (f *Repo) AdvertiseRefs(service string, w io.Writer) error {
	srv := server.NewServer(f.loader)

	var err error
	var ar *packp.AdvRefs

	switch service {
	case "git-receive-pack":
		var sess transport.ReceivePackSession
		if sess, err = srv.NewReceivePackSession(&transport.Endpoint{}, nil); err != nil {
			return err
		} else {
			ar, err = sess.AdvertisedReferencesContext(f.ctx)
		}
	case "git-upload-pack":
		var sess transport.UploadPackSession
		if sess, err = srv.NewUploadPackSession(&transport.Endpoint{}, nil); err != nil {
			return err
		} else {
			ar, err = sess.AdvertisedReferencesContext(f.ctx)
		}
	}

	if err != nil {
		return err
	}

	enc := pktline.NewEncoder(w)
	enc.Encodef("# service=%s\n", service)
	enc.Flush()

	return ar.Encode(w)
}

// ReceivePack
func (f *Repo) ReceivePack(r io.Reader, w io.Writer) error {
	srv := server.NewServer(f.loader)

	if sess, err := srv.NewReceivePackSession(&transport.Endpoint{}, nil); err != nil {
		return err
	} else {
		req := packp.NewReferenceUpdateRequest()
		if err := req.Decode(r); err != nil {
			return err
		}

		if status, _ := sess.ReceivePack(f.ctx, req); status != nil {
			return status.Encode(w)
		}

		return nil
	}
}

// UploadPack
func (f *Repo) UploadPack(r io.Reader, w io.Writer) error {
	srv := server.NewServer(f.loader)

	if sess, err := srv.NewUploadPackSession(&transport.Endpoint{}, nil); err != nil {
		return err
	} else {
		req := packp.NewUploadPackRequest()
		if err := req.Decode(r); err != nil {
			return err
		}

		if status, _ := sess.UploadPack(f.ctx, req); status != nil {
			return status.Encode(w)
		}

		return nil
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
