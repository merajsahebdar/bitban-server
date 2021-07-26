package ssh

import (
	"context"
	"net"
)

// srvContextKey
type srvContextKey struct{}

// netConnContextKey
type netConnContextKey struct{}

// newContext
func newContext(srv *Server, netConn net.Conn) (nextCtx context.Context, cancel context.CancelFunc) {
	nextCtx, cancel = context.WithCancel(context.Background())

	nextCtx = context.WithValue(nextCtx, srvContextKey{}, srv)
	nextCtx = context.WithValue(nextCtx, netConnContextKey{}, netConn)

	return nextCtx, cancel
}
