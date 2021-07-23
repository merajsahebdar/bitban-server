package common

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

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
			nextCtx = ContextWithDb(nextCtx, GetDbInstance())

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

// GetCookie
func GetCookie(ctx context.Context, cookieName string) (cookie *http.Cookie, err error) {
	ec := getEchoContext(ctx)
	if cookie, err = ec.Cookie(cookieName); err != nil {
		return nil, ErrMissingCookie
	}

	return cookie, nil
}
