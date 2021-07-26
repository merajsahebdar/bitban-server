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

package ssh

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
)

// CmdLine
type CmdLine string

// requestCmd
type requestCmd struct {
	Line CmdLine
}

// requestEnv
type requestEnv struct {
	Key   string
	Value string
}

// exitStatus
type exitResponse struct {
	Status uint
}

// newSession
func newSession(srv *Server, sshConn *gossh.ServerConn, ctx context.Context, ch gossh.NewChannel) {
	if curr, reqs, err := ch.Accept(); err != nil {
		srv.log.Error("got an error on channel accept", zap.Error(err))
	} else {
		go func(reqs <-chan *gossh.Request) {
			defer curr.Close()

			var cmd requestCmd
			var vars []requestEnv
			var isHandled bool

			for req := range reqs {
				switch req.Type {
				case "env":
					if isHandled {
						req.Reply(false, nil)
						continue
					}

					var env requestEnv
					gossh.Unmarshal(req.Payload, &env)
					vars = append(vars, env)
				case "exec":
					if isHandled {
						req.Reply(false, nil)
						continue
					}

					gossh.Unmarshal(req.Payload, &cmd)

					isHandled = true
					req.Reply(true, nil)

					go func() {
						srv.log.Info(fmt.Sprintf("%+v\n%+v\n", cmd, vars))
						exitCh(curr, 0)
					}()
				default:
					req.Reply(false, []byte(fmt.Sprintf("unsupported request: %s\n", req.Type)))
				}
			}
		}(reqs)
	}
}

// exitCh
func exitCh(ch gossh.Channel, status uint) {
	res := exitResponse{Status: status}
	ch.SendRequest("exit-status", false, gossh.Marshal(res))
	ch.Close()
}
