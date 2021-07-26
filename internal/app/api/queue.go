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

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.uber.org/fx"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/queue"
)

// QueueOpt
var QueueOpt = fx.Options(fx.Provide(newQueue), fx.Invoke(registerQueueLifecycle))

// newPoisonHandler
func newPoisonHandler() (mid message.HandlerMiddleware) {
	var err error

	if mid, err = middleware.PoisonQueue(
		queue.GetPublisherInstance(),
		"poisons",
	); err != nil {
		conf.Log.Fatal(err.Error())
	}

	return mid
}

// newQueue
func newQueue() *message.Router {
	var err error

	var router *message.Router
	if router, err = message.NewRouter(
		message.RouterConfig{},
		watermill.NopLogger{},
	); err != nil {
		conf.Log.Fatal(err.Error())
	}

	router.AddMiddleware(
		newPoisonHandler(),
	)

	return router
}

// registerQueueLifecycle
func registerQueueLifecycle(lc fx.Lifecycle, r *message.Router) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go func() {
				if err := r.Run(context.Background()); err != nil {
					conf.Log.Fatal(err.Error())
				}
			}()
			<-r.Running()

			conf.Log.Info("queue router is running...")

			return nil
		},
	})
}
