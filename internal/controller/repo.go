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
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"regeet.io/api/internal/facade"
)

// Repo
type Repo struct{}

// InfoRefs
func (c *Repo) InfoRefs(ec echo.Context) error {
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
func (c *Repo) ReceivePack(ec echo.Context) error {
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

	res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", transport.ReceivePackServiceName))
	res.Header().Set("Cache-Control", "no-cache")
	res.WriteHeader(200)

	if err := repo.ReceivePack(req.Body, res.Writer, false); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// UploadPack
func (c *Repo) UploadPack(ec echo.Context) error {
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

	res.Header().Set("Content-Type", fmt.Sprintf("application/x-%s-result", transport.ReceivePackServiceName))
	res.Header().Set("Cache-Control", "no-cache")
	res.WriteHeader(200)

	if err := repo.UploadPack(req.Body, res.Writer, false); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return nil
}

// RepoOpt
var RepoOpt = fx.Provide(newRepo)

// newRepo
func newRepo() *Repo {
	return &Repo{}
}
