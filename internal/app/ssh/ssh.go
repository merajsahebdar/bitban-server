package ssh

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	gossh "golang.org/x/crypto/ssh"
	"regeet.io/api/internal/conf"
)

// Ssh
type Ssh struct {
	config *gossh.ServerConfig
}

// ListenAndServe
func (s *Ssh) ListenAndServe(listener net.Listener, handler func() error) {
	for {
		if netConn, err := listener.Accept(); err != nil {
			conf.Log.Error("ssh: got an error on acception connection", zap.Error(err))
		} else {
			if sshConn, chans, reqs, err := gossh.NewServerConn(netConn, s.config); err != nil {
				if err == io.EOF {
					conf.Log.Error("ssh: handshaking was terminated")
				} else {
					conf.Log.Error("ssh: error on handshaking", zap.Error(err))
				}
			} else {
				go gossh.DiscardRequests(reqs)
				go s.resolveConnection(sshConn, chans)
			}
		}
	}
}

// resolveConnection
func (ssh *Ssh) resolveConnection(sshConn gossh.Conn, chans <-chan gossh.NewChannel) {
	for ch := range chans {
		if curr, reqs, err := ch.Accept(); err != nil {
			continue
		} else {
			go func(reqs <-chan *gossh.Request) {
				defer curr.Close()

				for req := range reqs {
					switch req.Type {
					case "env":
						conf.Log.Info("got an env request")
						continue
					case "exec":
						conf.Log.Info("got an exec request")
						return
					default:
						curr.Write([]byte(fmt.Sprintf("request not supported %s\n", req.Type)))
						return
					}
				}
			}(reqs)
		}
	}
}

// SshOpt
var SshOpt = fx.Provide(newSsh)

// newSsh
func newSsh() *Ssh {
	var err error

	var privatePEM []byte
	if privatePEM, err = base64.StdEncoding.DecodeString(conf.Cog.Ssh.Key.PrivateKey); err != nil {
		conf.Log.Fatal("failed to decode jwt private key")
	}

	var privateKey gossh.Signer
	if privateKey, err = gossh.ParsePrivateKeyWithPassphrase(
		privatePEM,
		[]byte(conf.Cog.Ssh.Key.Passphrase),
	); err != nil {
		conf.Log.Fatal("failed to parse jwt private key")
	}

	sshConfig := &gossh.ServerConfig{}
	sshConfig.NoClientAuth = true
	sshConfig.AddHostKey(privateKey)

	return &Ssh{
		config: sshConfig,
	}
}
