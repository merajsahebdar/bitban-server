package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.giteam.ir/giteam/internal/component"
	"go.giteam.ir/giteam/internal/fault"
	"go.giteam.ir/giteam/internal/util"
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
func isTokenExpired(claims *jwt.StandardClaims) bool {
	currUnix := time.Now().In(time.UTC).Unix()
	return claims.ExpiresAt <= currUnix
}

// GetContextRefreshTokenClaims
//
// Errors:
//   - fault.ErrInvalidJwtToken
// ErrorsRef:
//   - util.GetCookie
func GetContextRefreshTokenClaims(ctx context.Context) (*jwt.StandardClaims, error) {
	if cookie, err := util.GetCookie(ctx, refreshTokenCookie); err != nil {
		return nil, err
	} else {
		if claims, err := component.GetJwtInstance().VerifyToken(cookie.Value); err != nil || isTokenExpired(claims) {
			return nil, fault.ErrInvalidJwtToken
		} else {
			return claims, nil
		}
	}
}

// GetContextAccessTokenClaims
//
// Errors:
//   - fault.ErrMissingJwtToken in case of missing jwt token
//   - fault.ErrInvalidJwtToken in case of invalid or expired jwt token
func GetContextAccessTokenClaims(ctx context.Context) (*jwt.StandardClaims, error) {
	var token string

	authorization := util.GetHeader(ctx, echo.HeaderAuthorization)
	schemeLength := len(authHeaderScheme)
	if len(authorization) > schemeLength+1 && authorization[:schemeLength] == authHeaderScheme {
		token = authorization[schemeLength+1:]
	}

	if token == "" {
		return nil, fault.ErrMissingJwtToken
	}

	if claims, err := component.GetJwtInstance().VerifyToken(token); err != nil || isTokenExpired(claims) {
		return nil, fault.ErrInvalidJwtToken
	} else {
		return claims, nil
	}
}
