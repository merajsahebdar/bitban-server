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

	"github.com/go-git/go-git/v5/plumbing/transport"
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

	repo.AdvertiseRefs(ec.QueryParam("service"), res.Writer)

	return nil
}

// ReceivePack
func (c *Repo) ReceivePack(ctx context.Context) error {
	var repoName string
	var advRequest bool

	var r io.Reader
	var w io.Writer

	if ec, err := util.GetEchoContext(ctx); err != nil {
		if ch, err := ssh.GetContextCh(ctx); err != nil {
			// TODO: what we should do?
		} else {
			cmd := ssh.GetContextCmd(ctx)
			repoName = cmd.Args
			advRequest = true

			r = ioutil.NopCloser(ch)
			w = ch
		}
	} else {
		repoName = ec.Param("name")
		advRequest = false

		req := ec.Request()
		res := ec.Response()

		r = req.Body
		w = res.Writer

		res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", transport.ReceivePackServiceName))
		res.Header().Set("Cache-Control", "no-cache")
		res.WriteHeader(200)
	}

	var err error
	var repo *facade.Repo

	if repo, err = facade.GetRepoByName(
		ctx,
		strings.TrimSuffix(repoName, ".git"),
	); err != nil {
		// TODO:
		cfg.Log.Error("failed to initiate a repo facade", zap.Error(err))
		return nil
	}

	if err := repo.ReceivePack(r, w, advRequest); err != nil {
		// TODO: what is the best way to handle this error?
		cfg.Log.Error("got an error on precessing git request", zap.Error(err))
	}

	return nil
}

// UploadPack
func (c *Repo) UploadPack(ctx context.Context) error {
	var repoName string
	var advRequest bool

	var r io.Reader
	var w io.Writer

	if ec, err := util.GetEchoContext(ctx); err != nil {
		if ch, err := ssh.GetContextCh(ctx); err != nil {
			// TODO: what we should do?
		} else {
			cmd := ssh.GetContextCmd(ctx)
			repoName = cmd.Args
			advRequest = true

			r = ioutil.NopCloser(ch)
			w = ch
		}
	} else {
		repoName = ec.Param("name")
		advRequest = false

		req := ec.Request()
		res := ec.Response()

		r = req.Body
		w = res.Writer

		res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", transport.ReceivePackServiceName))
		res.Header().Set("Cache-Control", "no-cache")
		res.WriteHeader(200)
	}

	var err error
	var repo *facade.Repo

	if repo, err = facade.GetRepoByName(
		ctx,
		strings.TrimSuffix(repoName, ".git"),
	); err != nil {
		// TODO:
		cfg.Log.Error("failed to initiate a repo facade", zap.Error(err))
		return nil
	}

	if err := repo.UploadPack(r, w, advRequest); err != nil {
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
