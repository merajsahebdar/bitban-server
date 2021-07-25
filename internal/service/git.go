package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.giteam.ir/giteam/internal/conf"
	"go.giteam.ir/giteam/internal/facade"
	"go.uber.org/fx"
	"go.uber.org/zap"
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

// GitOpt
var GitOpt = fx.Invoke(registerGitHandlers)

// registerGitHandlers
func registerGitHandlers(ee *echo.Echo) {
	svc := &gitService{}

	//
	// Repository Handlers

	eg := ee.Group("/-/:name", func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conf.Log.Info("got a git client request", zap.String("path", c.Request().URL.Path), zap.String("query", c.Request().URL.RawQuery))
			return hf(c)
		}
	})
	eg.GET("/info/refs", svc.InfoRefs)
	eg.POST("/git-receive-pack", svc.ReceivePack)
	eg.POST("/git-upload-pack", svc.UploadPack)
}
