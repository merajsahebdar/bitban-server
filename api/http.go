package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.giteam.ir/giteam/internal/common"
	"go.giteam.ir/giteam/internal/conf"
	"go.giteam.ir/giteam/internal/resolver"
	"go.giteam.ir/giteam/internal/schema"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// graphQLQueryRoute
const graphQLQueryRoute = "/"

// newGraphQLQueryHandler Returns GraphQL query handler.
func newGraphQLQueryHandler(schemaConfig schema.Config) echo.HandlerFunc {
	// Query Handler
	queryHandler := handler.NewDefaultServer(schema.NewExecutableSchema(schemaConfig))

	// Panic Recover Handler
	queryHandler.SetRecoverFunc(func(ctx context.Context, mayErr interface{}) (userError error) {
		common.SetResponseStatus(ctx, http.StatusInternalServerError)

		fields := []zapcore.Field{}

		switch err := mayErr.(type) {
		case error:
			fields = append(fields, zap.Error(err))
		case string:
			fields = append(fields, zap.String("error", err))
		}

		conf.Log.Error("got a panic error when processing an api request", fields...)

		return resolver.InternalServerErrorFrom(nil)
	})

	// Enable tracing in development mode.
	if conf.CurrentEnv == conf.Dev {
		queryHandler.Use(apollotracing.Tracer{})
	}

	return func(ec echo.Context) error {
		queryHandler.ServeHTTP(ec.Response(), ec.Request())
		return nil
	}
}

// newGraphQLPlaygroundHandler Returns GraphQL playground handler.
func newGraphQLPlaygroundHandler() echo.HandlerFunc {
	playgroundHandler := playground.Handler("GraphQL Playground", graphQLQueryRoute)

	return func(ec echo.Context) error {
		playgroundHandler.ServeHTTP(ec.Response(), ec.Request())
		return nil
	}
}

// HttpOpt
var HttpOpt = fx.Options(fx.Provide(newHttp), fx.Invoke(registerHttpLifecycle))

// newHttp
func newHttp(schemaConfig schema.Config) *echo.Echo {
	e := echo.New()

	e.Use(common.ContextWrapper())
	e.Use(middleware.Recover())

	e.POST(graphQLQueryRoute, newGraphQLQueryHandler(schemaConfig))

	// Register playground just in development mode.
	if conf.CurrentEnv == conf.Dev {
		e.GET("/playground", newGraphQLPlaygroundHandler())
	}

	return e
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
