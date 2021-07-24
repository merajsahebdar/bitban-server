package api

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"go.giteam.ir/giteam/internal/conf"
	"go.giteam.ir/giteam/internal/queue"
	"go.uber.org/fx"
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
