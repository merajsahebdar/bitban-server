package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.giteam.ir/giteam/internal/common"
)

const (
	// refreshTokenCookie
	refreshTokenCookie = "refresh-token"

	// authHeaderScheme
	authHeaderScheme = "Bearer"
)

// SetRefreshTokenCookie
func SetRefreshTokenCookie(ctx context.Context, refreshToken string) {
	common.SetCookie(ctx, &http.Cookie{
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
	if cookie, err := common.GetCookie(ctx, refreshTokenCookie); err != nil {
		return nil, err
	} else {
		if claims, err := common.GetJwtInstance().VerifyToken(cookie.Value); err != nil || isTokenExpired(claims) {
			return nil, common.ErrInvalidJwtToken
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
	var token string

	authorization := common.GetHeader(ctx, echo.HeaderAuthorization)
	schemeLength := len(authHeaderScheme)
	if len(authorization) > schemeLength+1 && authorization[:schemeLength] == authHeaderScheme {
		token = authorization[schemeLength+1:]
	}

	if token == "" {
		return nil, common.ErrMissingJwtToken
	}

	if claims, err := common.GetJwtInstance().VerifyToken(token); err != nil || isTokenExpired(claims) {
		return nil, common.ErrInvalidJwtToken
	} else {
		return claims, nil
	}
}
