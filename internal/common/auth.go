package common

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// refreshTokenCookie
const refreshTokenCookie = "refresh-token"

// SetRefreshTokenCookie
func SetRefreshTokenCookie(ctx context.Context, refreshToken string) {
	SetCookie(ctx, &http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		HttpOnly: true,
	})
}

// GetContextRefreshToken
func GetContextRefreshTokenClaims(ctx context.Context) (*jwt.StandardClaims, error) {
	if cookie, err := GetCookie(ctx, refreshTokenCookie); err != nil {
		return nil, err
	} else {
		if claims, err := GetJwtInstance().VerifyToken(cookie.Value); err != nil {
			return nil, err
		} else {
			currTime := time.Now().In(time.UTC)
			if claims.ExpiresAt > currTime.Unix() {
				return nil, ErrInvalidJwtToken
			}

			return claims, nil
		}
	}
}
