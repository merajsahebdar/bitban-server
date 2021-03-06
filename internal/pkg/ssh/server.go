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
	"encoding/base64"
	"io"
	"net"

	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
	"bitban.io/server/internal/cfg"
)

// HandlerFunc
type HandlerFunc func(ctx context.Context) error

// Server
type Server struct {
	log     *sshLog
	cfgig  *gossh.ServerConfig
	handler map[string]HandlerFunc
}

// Use
func (srv *Server) Use(cmdLine string, handler HandlerFunc) {
	if _, ok := srv.handler[cmdLine]; ok {
		srv.log.Fatal("not allowed to register more than one handler for each command")
	} else {
		srv.handler[cmdLine] = handler
	}
}

// ListenAndServe
func (srv *Server) ListenAndServe(listener net.Listener) {
	for {
		if netConn, err := listener.Accept(); err != nil {
			srv.log.Error("ssh: got an error on acception connection", zap.Error(err))
		} else {
			go func() {
				ctx, cancel := newContext(srv, netConn)

				defer func() {
					netConn.Close()
					cancel()
				}()

				if sshConn, chans, reqs, err := gossh.NewServerConn(netConn, srv.cfgig); err != nil {
					if err == io.EOF {
						srv.log.Error("handshaking was terminated")
					} else {
						srv.log.Error("error on handshaking", zap.Error(err))
					}
				} else {
					go gossh.DiscardRequests(reqs)
					for ch := range chans {
						newSession(srv, sshConn, ctx, ch)
					}
				}
			}()
		}
	}
}

// NewServer
func NewServer(logScope string) *Server {
	var err error

	log := &sshLog{
		cfg.Log.With(zap.String("scope", logScope)),
	}

	var privatePEM []byte
	if privatePEM, err = base64.StdEncoding.DecodeString(cfg.Cog.Ssh.Key.PrivateKey); err != nil {
		log.Fatal("failed to decode private key")
	}

	var privateKey gossh.Signer
	if privateKey, err = gossh.ParsePrivateKeyWithPassphrase(
		privatePEM,
		[]byte(cfg.Cog.Ssh.Key.Passphrase),
	); err != nil {
		log.Fatal("failed to parse private key")
	}

	sshConfig := &gossh.ServerConfig{}
	sshConfig.NoClientAuth = true
	sshConfig.AddHostKey(privateKey)

	return &Server{
		log:     log,
		cfgig:  sshConfig,
		handler: make(map[string]HandlerFunc),
	}
}
