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
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
)

// RequestCmd
type RequestCmd struct {
	Line string
	Name string
	Args string
}

// requestEnv
type requestEnv struct {
	Key   string
	Value string
}

// session
type session struct {
	*sshLog
	sync.Mutex
	gossh.Channel
	handler map[string]HandlerFunc
	reqs    <-chan *gossh.Request
	ctx     context.Context
	exited  bool
}

// exit
func (sess *session) exit(status []byte) error {
	sess.Lock()
	defer sess.Unlock()

	if sess.exited {
		return errors.New("session was exited before")
	}

	sess.exited = true

	if _, err := sess.SendRequest("exit-status", false, status); err != nil {
		return err
	}

	return sess.Close()
}

// track
func (sess *session) track() {
	go func(reqs <-chan *gossh.Request) {
		defer sess.Close()

		var cmd RequestCmd
		var envs []requestEnv
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
				envs = append(envs, env)
				req.Reply(true, nil)
			case "exec":
				if isHandled {
					req.Reply(false, nil)
					continue
				}

				parseExecCommand(req.Payload, &cmd)

				isHandled = true
				req.Reply(true, nil)

				go func() {
					sess.Info("received a new request", zap.String("request", cmd.Name))
					if handler, ok := sess.handler[cmd.Name]; ok {
						var nextCtx context.Context
						nextCtx = withContextCmd(sess.ctx, cmd)
						nextCtx = withContextEnvs(nextCtx, envs)
						nextCtx = withContextCh(nextCtx, sess.Channel)
						handler(nextCtx)
						sess.exit([]byte{0, 0, 0, 0})
					} else {
						sess.Error(
							"no handler was found for the requested command",
							zap.String("cmd", cmd.Name),
						)
					}
				}()
			default:
				req.Reply(false, []byte(fmt.Sprintf("unsupported request: %s\n", req.Type)))
			}
		}
	}(sess.reqs)
}

// newSession
func newSession(srv *Server, sshConn *gossh.ServerConn, ctx context.Context, ch gossh.NewChannel) {
	if curr, reqs, err := ch.Accept(); err != nil {
		srv.log.Error("got an error on channel accept", zap.Error(err))
	} else {
		sess := &session{
			Mutex:   sync.Mutex{},
			Channel: curr,
			reqs:    reqs,
			sshLog:  srv.log,
			ctx:     ctx,
			handler: srv.handler,
		}

		sess.track()
	}
}
