package api

import (
	"context"
	"fmt"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/schema"
	"regeet.io/api/internal/util"
)

// HttpOpt
var HttpOpt = fx.Options(fx.Provide(newHttp), fx.Invoke(registerHttpLifecycle))

// newHttp
func newHttp(schemaConfig schema.Config) *echo.Echo {
	ee := echo.New()
	ee.Use(util.ContextWrapper())
	ee.Use(middleware.Recover())

	return ee
}

// registerHttpLifecycle
func registerHttpLifecycle(lifecycle fx.Lifecycle, e *echo.Echo) {
	e.HideBanner = true
	e.HidePort = true

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			addr := fmt.Sprintf("%s:%d", conf.Cog.App.Host, conf.Cog.App.Port)

			if e.Listener, err = net.Listen("tcp", addr); err != nil {
				conf.Log.Fatal("cannot start the http listener", zap.Error(err))
			}

			conf.Log.Info("ready to respond http requests...", zap.String("addr", addr))

			go func() {
				if err := e.Start(addr); err != nil {
					conf.Log.Fatal("cannot start the http server", zap.Error(err))
				}
			}()

			return nil
		},
	})
}
