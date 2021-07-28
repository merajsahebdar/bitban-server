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
	"net"

	gossh "golang.org/x/crypto/ssh"
)

// srvContextKey
type srvContextKey struct{}

// netConnContextKey
type netConnContextKey struct{}

// envsContextKey
type envsContextKey struct{}

// cmdContextKey
type cmdContextKey struct{}

// chContextKey
type chContextKey struct{}

// newContext
func newContext(srv *Server, netConn net.Conn) (nextCtx context.Context, cancel context.CancelFunc) {
	nextCtx, cancel = context.WithCancel(context.Background())

	nextCtx = context.WithValue(nextCtx, srvContextKey{}, srv)
	nextCtx = context.WithValue(nextCtx, netConnContextKey{}, netConn)

	return nextCtx, cancel
}

// GetContextConn
func GetContextConn(ctx context.Context) net.Conn {
	return ctx.Value(netConnContextKey{}).(net.Conn)
}

// withContextEnvs
func withContextEnvs(ctx context.Context, envs []requestEnv) context.Context {
	return context.WithValue(
		ctx,
		envsContextKey{},
		envs,
	)
}

// GetContextEnvs
func GetContextEnvs(ctx context.Context) []requestEnv {
	return ctx.Value(envsContextKey{}).([]requestEnv)
}

// withContextCh
func withContextCh(ctx context.Context, ch gossh.Channel) context.Context {
	return context.WithValue(
		ctx,
		chContextKey{},
		ch,
	)
}

// GetContextCh
func GetContextCh(ctx context.Context) (gossh.Channel, error) {
	if ch, ok := ctx.Value(chContextKey{}).(gossh.Channel); ok {
		return ch, nil
	} else {
		return nil, errors.New("no ssh channel")
	}
}

// MustGetContextCh
func MustGetContextCh(ctx context.Context) gossh.Channel {
	if ch, err := GetContextCh(ctx); err != nil {
		panic(err)
	} else {
		return ch
	}
}

// withContextCmd
func withContextCmd(ctx context.Context, cmd RequestCmd) context.Context {
	return context.WithValue(
		ctx,
		cmdContextKey{},
		cmd,
	)
}

// GetContextCmd
func GetContextCmd(ctx context.Context) RequestCmd {
	return ctx.Value(cmdContextKey{}).(RequestCmd)
}
