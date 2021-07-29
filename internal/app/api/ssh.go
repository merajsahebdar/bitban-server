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

package api

import (
	"context"
	"net"

	"go.uber.org/fx"
	"regeet.io/api/internal/app/controller"
	"regeet.io/api/internal/cfg"
	"regeet.io/api/internal/pkg/facade"
	"regeet.io/api/internal/pkg/ssh"
)

// SshOpt
var SshOpt = fx.Invoke(registerSshServerLifecycle)

// registerSshServerLifecycle
func registerSshServerLifecycle(lc fx.Lifecycle, repoController *controller.Repo) {
	srv := ssh.NewServer("api.ssh")

	srv.Use(facade.GitReceivePack, repoController.ServePack)
	srv.Use(facade.GitUploadPack, repoController.ServePack)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			var err error

			var listener net.Listener
			if listener, err = net.Listen("tcp", ":8022"); err != nil {
				cfg.Log.Fatal("service: cannot start the ssh listener")
			}

			go func() {
				srv.ListenAndServe(listener)
			}()

			return nil
		},
	})
}
