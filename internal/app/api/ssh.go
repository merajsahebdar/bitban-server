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
	"io/ioutil"
	"net"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"go.uber.org/fx"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/facade"
	"regeet.io/api/internal/pkg/ssh"
)

// SshOpt
var SshOpt = fx.Invoke(registerSshServerLifecycle)

// registerSshServerLifecycle
func registerSshServerLifecycle(lc fx.Lifecycle) {
	srv := ssh.NewServer("api.ssh")

	srv.Use(transport.UploadPackServiceName, func(ctx context.Context, args string) error {
		ch := ssh.GetContextCh(ctx)

		if repo, err := facade.GetRepoByName(
			ctx,
			strings.TrimSuffix(args, ".git"),
		); err != nil {
			return err
		} else {
			if err := repo.UploadPack(ioutil.NopCloser(ch), ch, true); err != nil {
				return err
			}

			return nil
		}
	})

	srv.Use(transport.ReceivePackServiceName, func(ctx context.Context, args string) error {
		ch := ssh.GetContextCh(ctx)

		if repo, err := facade.GetRepoByName(
			ctx,
			strings.TrimSuffix(args, ".git"),
		); err != nil {
			return err
		} else {
			if err := repo.ReceivePack(ioutil.NopCloser(ch), ch, true); err != nil {
				return err
			}

			return nil
		}
	})

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
