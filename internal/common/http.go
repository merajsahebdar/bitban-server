package common

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.giteam.ir/giteam/internal/dto"
	"go.giteam.ir/giteam/internal/orm"
)

// authScheme
const authScheme = "Bearer"

// AuthCookie
const AuthCookie = "refresh-token"

// authorizedUserKey
type authorizedUserKey struct{}

// echoContextKey Key to access the echo context
type echoContextKey struct{}

// echoContext
type echoContext struct {
	echo.Context
	ctx context.Context
}

// ContextWrapper Wraps Echo context to keep it for future uses.
func ContextWrapper() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ec echo.Context) error {
			var nextCtx context.Context
			nextCtx = context.WithValue(ec.Request().Context(), echoContextKey{}, ec)
			nextCtx = ContextWithDB(nextCtx, GetDbInstance())

			ec.SetRequest(ec.Request().WithContext(nextCtx))

			return next(echoContext{
				Context: ec,
				ctx:     nextCtx,
			})
		}
	}
}

// getEchoContext
func getEchoContext(ctx context.Context) echo.Context {
	return ctx.Value(echoContextKey{}).(echo.Context)
}

// SetResponseStatus
func SetResponseStatus(ctx context.Context, status int) {
	res := getEchoContext(ctx).Response()
	res.Status = status
}

// SetCookie
func SetCookie(ctx context.Context, cookie *http.Cookie) {
	ec := getEchoContext(ctx)
	ec.SetCookie(cookie)
}

// GetAuthorizedUser
func GetAuthorizedUser(ctx context.Context) *orm.User {
	if user, ok := ctx.Value(authorizedUserKey{}).(*orm.User); ok {
		return user
	}

	return nil
}

// AuthorizeUser
func AuthorizeUser(ctx context.Context, subLookup dto.PermissionSubLookup) (*orm.User, error) {
	if user := GetAuthorizedUser(ctx); user != nil {
		return user, nil
	}

	ec := getEchoContext(ctx)

	var err error
	var token string

	switch subLookup {
	case dto.PermissionSubHeaderLookup:
		authorization := ec.Request().Header.Get(echo.HeaderAuthorization)
		schemeLength := len(authScheme)
		if len(authorization) > schemeLength+1 && authorization[:schemeLength] == authScheme {
			token = authorization[schemeLength+1:]
			break
		}
	case dto.PermissionSubCookieLookup:
		var cookie *http.Cookie
		if cookie, err = ec.Cookie(AuthCookie); err != nil {
			return nil, ErrMissingJwtToken
		}

		token = cookie.Value
	}

	if token == "" {
		return nil, ErrMissingJwtToken
	}

	var claims *jwt.StandardClaims
	if claims, err = GetJwtInstance().VerifyToken(token); err != nil {
		return nil, ErrInvalidJwtToken
	}

	currentTime := time.Now().Unix()
	if claims.ExpiresAt < currentTime || claims.NotBefore > currentTime {
		return nil, ErrInvalidJwtToken
	}

	var id int64
	if id, err = strconv.ParseInt(claims.Subject, 10, 64); err != nil {
		return nil, ErrInvalidJwtToken
	}

	var user *orm.User
	if user, err = orm.FindUser(ctx, GetDbInstance(), id); err != nil {
		return nil, ErrInvalidJwtToken
	}

	return user, nil
}

// WithAuthorizedUser
func WithAuthorizedUser(ctx context.Context, user *orm.User) context.Context {
	return context.WithValue(ctx, authorizedUserKey{}, user)
}
