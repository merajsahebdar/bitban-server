/*
 * Copyright 2021 Meraj Sahebdar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/resolver"
	"regeet.io/api/internal/schema"
	"regeet.io/api/internal/util"
)

// EchoOpt
var EchoOpt = fx.Invoke(registerEchoLifecycle)

// registerEchoLifecycle
func registerEchoLifecycle(lc fx.Lifecycle, schemaConfig schema.Config) {
	ee := echo.New()
	ee.Use(util.ContextWrapper())
	ee.Use(middleware.Recover())

	//
	// Register Git

	svc := &gitService{}

	eg := ee.Group("/-/:name", func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conf.Log.Info("got a git client request", zap.String("path", c.Request().URL.Path), zap.String("query", c.Request().URL.RawQuery))
			return hf(c)
		}
	})
	eg.GET("/info/refs", svc.InfoRefs)
	eg.POST("/git-receive-pack", svc.ReceivePack)
	eg.POST("/git-upload-pack", svc.UploadPack)

	//
	// Register GraphQL

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

		conf.Log.Error("got a panic error when processing a graphql request", fields...)

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

	ee.HideBanner = true
	ee.HidePort = true

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			addr := fmt.Sprintf("%s:%d", conf.Cog.App.Host, conf.Cog.App.Port)

			if ee.Listener, err = net.Listen("tcp", addr); err != nil {
				conf.Log.Fatal("cannot start the http listener", zap.Error(err))
			}

			conf.Log.Info("ready to respond http requests...", zap.String("addr", addr))

			go func() {
				if err := ee.Start(addr); err != nil {
					conf.Log.Fatal("cannot start the http server", zap.Error(err))
				}
			}()

			return nil
		},
	})
}
