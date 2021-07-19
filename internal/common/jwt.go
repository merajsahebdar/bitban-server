package common

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"sync"

	"github.com/dgrijalva/jwt-go"
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

			if publicPEM, err = base64.StdEncoding.DecodeString(Cog.Jwt.PublicKey); err != nil {
				Log.Fatal("failed to decode jwt public key")
			}

			if privatePEM, err = base64.StdEncoding.DecodeString(Cog.Jwt.PrivateKey); err != nil {
				Log.Fatal("failed to decode jwt private key")
			}

			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(
				publicPEM,
			)
			if err != nil {
				Log.Fatal("failed to parse jwt public key")
			}

			privateKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(
				privatePEM,
				Cog.Jwt.Passphrase,
			)
			if err != nil {
				Log.Fatal("failed to parse jwt private key")
			}

			jwtInstance = &Jwt{
				publicKey:  publicKey,
				privateKey: privateKey,
			}
		}
	}

	return jwtInstance
}
