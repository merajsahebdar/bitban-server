package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"regeet.io/api/internal/app/ssh"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/facade"
)

// gitService
type gitService struct{}

// InfoRefs
func (s *gitService) InfoRefs(ec echo.Context) error {
	var err error
	var repo *facade.Repo

	req := ec.Request()
	res := ec.Response()

	if repo, err = facade.GetRepoByName(
		req.Context(),
		ec.Param("name"),
	); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-advertisement", ec.QueryParam("service")))
	res.Header().Set("Cache-Control", "no-cache")
	res.WriteHeader(200)

	repo.AdvertiseRefs(ec.QueryParam("service"), res.Writer)

	return nil
}

// ReceivePack
func (s *gitService) ReceivePack(ec echo.Context) error {
	var err error
	var repo *facade.Repo

	req := ec.Request()
	res := ec.Response()

	if repo, err = facade.GetRepoByName(
		req.Context(),
		ec.Param("name"),
	); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", "git-receive-pack"))
	res.Header().Set("Cache-Control", "no-cache")
	res.WriteHeader(200)

	if err := repo.ReceivePack(req.Body, res.Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// UploadPack
func (s *gitService) UploadPack(ec echo.Context) error {
	var err error
	var repo *facade.Repo

	req := ec.Request()
	res := ec.Response()

	if repo, err = facade.GetRepoByName(
		req.Context(),
		ec.Param("name"),
	); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", "git-upload-pack"))
	res.Header().Set("Cache-Control", "no-cache")
	res.WriteHeader(200)

	if err := repo.UploadPack(req.Body, res.Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// SshOpt
var SshOpt = fx.Invoke(registerSshServerLifecycle)

// registerSshServerLifecycle
func registerSshServerLifecycle(lc fx.Lifecycle, srv *ssh.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error

			var listener net.Listener
			if listener, err = net.Listen("tcp", ":8022"); err != nil {
				conf.Log.Fatal("service: cannot start the ssh listener")
			}

			go func() {
				srv.ListenAndServe(listener)
			}()

			return nil
		},
	})
}
