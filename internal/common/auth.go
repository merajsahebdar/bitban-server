package common

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	// refreshTokenCookie
	refreshTokenCookie = "refresh-token"

	// authHeaderScheme
	authHeaderScheme = "Bearer"
)

// SetRefreshTokenCookie
func SetRefreshTokenCookie(ctx context.Context, refreshToken string) {
	SetCookie(ctx, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		HttpOnly: true,
	})
}

// isTokenExpired
func isTokenExpired(claims *jwt.StandardClaims) bool {
	currUnix := time.Now().In(time.UTC).Unix()
	return claims.ExpiresAt <= currUnix
}

// GetContextRefreshTokenClaims
//
// Errors:
//   - common.ErrInvalidJwtToken
// ErrorsRef:
//   - common.GetCookie
func GetContextRefreshTokenClaims(ctx context.Context) (*jwt.StandardClaims, error) {
	if cookie, err := GetCookie(ctx, refreshTokenCookie); err != nil {
		return nil, err
	} else {
		if claims, err := GetJwtInstance().VerifyToken(cookie.Value); err != nil || isTokenExpired(claims) {
			return nil, ErrInvalidJwtToken
		} else {
			return claims, nil
		}
	}
}

// GetContextAccessTokenClaims
//
// Errors:
//   - common.ErrMissingJwtToken in case of missing jwt token
//   - common.ErrInvalidJwtToken in case of invalid or expired jwt token
func GetContextAccessTokenClaims(ctx context.Context) (*jwt.StandardClaims, error) {
	ec := getEchoContext(ctx)

	var token string

	authorization := ec.Request().Header.Get(echo.HeaderAuthorization)
	schemeLength := len(authHeaderScheme)
	if len(authorization) > schemeLength+1 && authorization[:schemeLength] == authHeaderScheme {
		token = authorization[schemeLength+1:]
	}

	if token == "" {
		return nil, ErrMissingJwtToken
	}

	if claims, err := GetJwtInstance().VerifyToken(token); err != nil || isTokenExpired(claims) {
		return nil, ErrInvalidJwtToken
	} else {
		return claims, nil
	}
}
