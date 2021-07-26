package ssh

import (
	"context"
	"encoding/base64"
	"io"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
	"regeet.io/api/internal/conf"
)

// Handler
type Handler func(ctx context.Context) error

// Server
type Server struct {
	log     *sshLog
	config  *gossh.ServerConfig
	handler map[CmdLine]Handler
}

// Use
func (srv *Server) Use(cmdLine CmdLine, handler Handler) {
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

				if sshConn, chans, reqs, err := gossh.NewServerConn(netConn, srv.config); err != nil {
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

// ServerOpt
var ServerOpt = fx.Provide(newLog, newServer)

// newServer
func newServer(log *sshLog) *Server {
	var err error

	var privatePEM []byte
	if privatePEM, err = base64.StdEncoding.DecodeString(conf.Cog.Ssh.Key.PrivateKey); err != nil {
		log.Fatal("failed to decode private key")
	}

	var privateKey gossh.Signer
	if privateKey, err = gossh.ParsePrivateKeyWithPassphrase(
		privatePEM,
		[]byte(conf.Cog.Ssh.Key.Passphrase),
	); err != nil {
		log.Fatal("failed to parse private key")
	}

	sshConfig := &gossh.ServerConfig{}
	sshConfig.NoClientAuth = true
	sshConfig.AddHostKey(privateKey)

	return &Server{
		log:     log,
		config:  sshConfig,
		handler: make(map[CmdLine]Handler),
	}
}
