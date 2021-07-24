package util

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.giteam.ir/giteam/internal/fault"
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
			nextCtx := context.WithValue(ec.Request().Context(), echoContextKey{}, ec)
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

// GetHeader
func GetHeader(ctx context.Context, key string) string {
	ec := getEchoContext(ctx)
	return ec.Request().Header.Get(key)
}

// GetCookie
//
// Errors:
//   - common.ErrMissingCookie in case of not founding the asked cookie
func GetCookie(ctx context.Context, cookieName string) (cookie *http.Cookie, err error) {
	ec := getEchoContext(ctx)
	if cookie, err = ec.Cookie(cookieName); err != nil {
		return nil, fault.ErrMissingCookie
	}

	return cookie, nil
}
