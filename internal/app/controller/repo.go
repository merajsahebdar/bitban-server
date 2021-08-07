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
	"regeet.io/api/internal/pkg/dto"
	"regeet.io/api/internal/pkg/facade"
	"regeet.io/api/internal/pkg/fault"
	"regeet.io/api/internal/pkg/orm/entity"
	"regeet.io/api/internal/pkg/ssh"
	"regeet.io/api/internal/pkg/util"
	"regeet.io/api/internal/pkg/validate"
)

// Repo
type Repo struct{}

// hasAccessByHttp
func (c *Repo) hasAccessByHttp(ctx context.Context, domain string, repo int64) bool {
	ec := util.MustGetEchoContext(ctx)
	req := ec.Request()
	res := ec.Response()

	authHeader := req.Header.Get(echo.HeaderAuthorization)
	if authHeader == "" {
		res.Header().Set(echo.HeaderWWWAuthenticate, `Basic realm="Restricted"`)
		return false
	}

	if identifier, password, ok := req.BasicAuth(); ok {
		if account, err := facade.GetAccountByPassword(ctx, dto.SignInInput{
			Identifier: identifier,
			Password:   password,
		}); err != nil {
			return false
		} else {
			if err := account.CheckPermissionIn(
				domain,
				fmt.Sprintf("/repositories/%d", repo),
				"git-serve",
			); err != nil {
				return false
			}
		}
	}

	return true
}

// InfoRefs
func (c *Repo) InfoRefs(ctx context.Context) error {
	ec := util.MustGetEchoContext(ctx)
	req := ec.Request()
	res := ec.Response()

	domainAddress := ec.Param("domain")
	repoAddress := ec.Param("repo")

	var err error
	var repo *facade.Repo

	if repo, err = facade.GetRepoByAddress(
		req.Context(),
		domainAddress,
		repoAddress,
	); err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if !c.hasAccessByHttp(ctx, domainAddress, repo.GetID()) {
		return echo.NewHTTPError(http.StatusForbidden)
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
	var domainAddress string
	var repoAddress string
	var isSsh bool

	var r io.Reader
	var w io.Writer

	if ec, err := util.GetEchoContext(ctx); err != nil {
		if ch, err := ssh.GetContextCh(ctx); err != nil {
			// TODO: what we should do?
		} else {
			cmd := ssh.GetContextCmd(ctx)

			service = cmd.Name

			paths := strings.SplitN(cmd.Args, "/", 2)
			domainAddress = paths[0]
			repoAddress = paths[1]

			isSsh = true

			r = ioutil.NopCloser(ch)
			w = ch
		}
	} else {
		service = ec.Param("service")
		domainAddress = ec.Param("domain")
		repoAddress = ec.Param("repo")
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

	if repo, err = facade.GetRepoByAddress(
		ctx,
		domainAddress,
		strings.TrimSuffix(repoAddress, ".git"),
	); err != nil {
		// TODO:
		cfg.Log.Error("failed to initiate a repo facade", zap.Error(err))
		return nil
	}

	if isSsh {
		// TODO: implement ssh-key auth
	} else {
		if !c.hasAccessByHttp(ctx, domainAddress, repo.GetID()) {
			return echo.NewHTTPError(http.StatusForbidden)
		}
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

// CreateRepository
//
// Errors:
//   - fault.ErrUnauthenticated if the request is not authorized
//   - fault.UserInputError if the provided input is invalid
// ErrorsRef:
//   - facade.CreateRepoByAddress
func (c *Repo) CreateRepository(ctx context.Context, input dto.CreateRepositoryInput) (*entity.Repository, error) {
	if currAccount, err := facade.GetAccountByAccessToken(ctx); err != nil {
		return nil, fault.ErrUnauthenticated
	} else {
		if err := validate.
			GetValidateInstance().
			Struct(input); err != nil {
			return nil, fault.UserInputErrorFrom(err)
		}

		if repo, err := facade.CreateRepoByAddress(ctx, currAccount.GetDomain().Address, input.Address); err != nil {
			return nil, err
		} else {
			return repo.GetEntity(), nil
		}
	}
}

// RepoOpt
var RepoOpt = fx.Provide(newRepo)

// newRepo
func newRepo() *Repo {
	return &Repo{}
}
