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

package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"regeet.io/api/internal/pkg/jwt"
	"regeet.io/api/internal/pkg/util"
)

var (
	// ErrMissingJwtToken
	ErrMissingJwtToken = errors.New("the jwt token is missing or malformed")

	// ErrInvalidJwtToken
	ErrInvalidJwtToken = errors.New("the jwt token is invalid or expired")
)

const (
	// refreshTokenCookie
	refreshTokenCookie = "refresh-token"

	// authHeaderScheme
	authHeaderScheme = "Bearer"
)

// SetRefreshTokenCookie
func SetRefreshTokenCookie(ctx context.Context, refreshToken string) {
	util.SetCookie(ctx, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		HttpOnly: true,
	})
}

// isTokenExpired
func isTokenExpired(claims *gojwt.StandardClaims) bool {
	currUnix := time.Now().In(time.UTC).Unix()
	return claims.ExpiresAt <= currUnix
}

// GetContextRefreshTokenClaims
//
// Errors:
//   - auth.ErrInvalidJwtToken
// ErrorsRef:
//   - util.GetCookie
func GetContextRefreshTokenClaims(ctx context.Context) (*gojwt.StandardClaims, error) {
	if cookie, err := util.GetCookie(ctx, refreshTokenCookie); err != nil {
		return nil, err
	} else {
		if claims, err := jwt.GetJwtInstance().VerifyToken(cookie.Value); err != nil || isTokenExpired(claims) {
			return nil, ErrInvalidJwtToken
		} else {
			return claims, nil
		}
	}
}

// GetContextAccessTokenClaims
//
// Errors:
//   - auth.ErrMissingJwtToken in case of missing jwt token
//   - auth.ErrInvalidJwtToken in case of invalid or expired jwt token
func GetContextAccessTokenClaims(ctx context.Context) (*gojwt.StandardClaims, error) {
	var token string

	authorization := util.GetHeader(ctx, echo.HeaderAuthorization)
	schemeLength := len(authHeaderScheme)
	if len(authorization) > schemeLength+1 && authorization[:schemeLength] == authHeaderScheme {
		token = authorization[schemeLength+1:]
	}

	if token == "" {
		return nil, ErrMissingJwtToken
	}

	if claims, err := jwt.GetJwtInstance().VerifyToken(token); err != nil || isTokenExpired(claims) {
		return nil, ErrInvalidJwtToken
	} else {
		return claims, nil
	}
}
