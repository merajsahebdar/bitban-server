package service

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/resolver"
	"regeet.io/api/internal/schema"
	"regeet.io/api/internal/util"
)

// GraphQLOpt
var GraphQLOpt = fx.Invoke(registerGraphQLHandlers)

// registerGraphQLHandlers
func registerGraphQLHandlers(ee *echo.Echo, schemaConfig schema.Config) {
	// Query Handler
	queryHandler := handler.NewDefaultServer(schema.NewExecutableSchema(schemaConfig))

	// Panic Recover Handler
	queryHandler.SetRecoverFunc(func(ctx context.Context, mayErr interface{}) (userError error) {
		util.SetResponseStatus(ctx, http.StatusInternalServerError)

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

	ee.POST("/api", func(ec echo.Context) error {
		queryHandler.ServeHTTP(ec.Response(), ec.Request())
		return nil
	})

	// Register playground just in development mode.
	if conf.CurrentEnv == conf.Dev {
		playgroundHandler := playground.Handler("GraphQL Playground", "/api")

		ee.GET("/api/playground", func(ec echo.Context) error {
			playgroundHandler.ServeHTTP(ec.Response(), ec.Request())
			return nil
		})
	}
}
