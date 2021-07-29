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

package controller

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"regeet.io/api/internal/cfg"
	"regeet.io/api/internal/pkg/facade"
	"regeet.io/api/internal/pkg/ssh"
	"regeet.io/api/internal/pkg/util"
)

// Repo
type Repo struct{}

// InfoRefs
func (c *Repo) InfoRefs(ctx context.Context) error {
	ec := util.MustGetEchoContext(ctx)

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

	repo.AdvertiseRefs(res.Writer, ec.QueryParam("service"))

	return nil
}

// ServePack
func (c *Repo) ServePack(ctx context.Context) error {
	var service string
	var name string
	var isSsh bool

	var r io.Reader
	var w io.Writer

	if ec, err := util.GetEchoContext(ctx); err != nil {
		if ch, err := ssh.GetContextCh(ctx); err != nil {
			// TODO: what we should do?
		} else {
			cmd := ssh.GetContextCmd(ctx)

			service = cmd.Name
			name = cmd.Args
			isSsh = true

			r = ioutil.NopCloser(ch)
			w = ch
		}
	} else {
		service = ec.Param("service")
		name = ec.Param("name")
		isSsh = false

		req := ec.Request()
		res := ec.Response()

		r = req.Body
		w = res.Writer

		res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", service))
		res.Header().Set("Cache-Control", "no-cache")
		res.WriteHeader(200)
	}

	var err error
	var repo *facade.Repo

	if repo, err = facade.GetRepoByName(
		ctx,
		strings.TrimSuffix(name, ".git"),
	); err != nil {
		// TODO:
		cfg.Log.Error("failed to initiate a repo facade", zap.Error(err))
		return nil
	}

	if err := repo.ServePack(&facade.ServerPackConfig{
		R:       r,
		W:       w,
		Service: service,
		IsSsh:   isSsh,
	}); err != nil {
		// TODO: what is the best way to handle this error?
		cfg.Log.Error("got an error on precessing git request", zap.Error(err))
	}

	return nil
}

// RepoOpt
var RepoOpt = fx.Provide(newRepo)

// newRepo
func newRepo() *Repo {
	return &Repo{}
}
