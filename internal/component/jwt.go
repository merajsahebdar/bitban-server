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

package component

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"regeet.io/api/internal/conf"
)

// Jwt
type Jwt struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// VerifyToken
func (j *Jwt) VerifyToken(token string) (*jwt.StandardClaims, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&jwt.StandardClaims{},
		func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodRSA)
			if !ok {
				return nil, errors.New("unexpected token signing method")
			}

			return j.publicKey, nil
		})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	c, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return c, nil
}

// SignToken
func (j *Jwt) SignToken(c *jwt.StandardClaims) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims = c
	return t.SignedString(j.privateKey)
}

// jwtComponentLock
var jwtComponentLock = &sync.Mutex{}

// jwtInstance
var jwtInstance *Jwt

// GetJwtInstance
func GetJwtInstance() *Jwt {
	if jwtInstance == nil {
		jwtComponentLock.Lock()
		defer jwtComponentLock.Unlock()

		if jwtInstance == nil {
			var err error
			var publicPEM, privatePEM []byte

			if publicPEM, err = base64.StdEncoding.DecodeString(conf.Cog.Jwt.PublicKey); err != nil {
				conf.Log.Fatal("failed to decode jwt public key")
			}

			if privatePEM, err = base64.StdEncoding.DecodeString(conf.Cog.Jwt.PrivateKey); err != nil {
				conf.Log.Fatal("failed to decode jwt private key")
			}

			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(
				publicPEM,
			)
			if err != nil {
				conf.Log.Fatal("failed to parse jwt public key")
			}

			privateKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(
				privatePEM,
				conf.Cog.Jwt.Passphrase,
			)
			if err != nil {
				conf.Log.Fatal("failed to parse jwt private key")
			}

			jwtInstance = &Jwt{
				publicKey:  publicKey,
				privateKey: privateKey,
			}
		}
	}

	return jwtInstance
}
