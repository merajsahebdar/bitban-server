package queue

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"go.giteam.ir/giteam/internal/conf"
)

// Queue
type Queue struct {
	subscriber *amqp.Subscriber
	publisher  *amqp.Publisher
}

// queueLock
var queueLock = &sync.Mutex{}

// queueInstance
var queueInstance *Queue

// getQueueInstance
func getQueueInstance() *Queue {
	if queueInstance == nil {
		queueLock.Lock()
		defer queueLock.Unlock()

		if queueInstance == nil {
			config := amqp.NewDurableQueueConfig(conf.Cog.Amqp.Uri)

			var err error

			var subscriber *amqp.Subscriber
			if subscriber, err = amqp.NewSubscriber(
				config,
				watermill.NopLogger{},
			); err != nil {
				conf.Log.Fatal(err.Error())
			}

			var publisher *amqp.Publisher
			if publisher, err = amqp.NewPublisher(
				config,
				watermill.NopLogger{},
			); err != nil {
				conf.Log.Fatal(err.Error())
			}

			queueInstance = &Queue{
				subscriber,
				publisher,
			}
		}
	}

	return queueInstance
}

// GetPublisherInstance
func GetPublisherInstance() *amqp.Publisher {
	return getQueueInstance().publisher
}

// GetSubscriberInstance
func GetSubscriberInstance() *amqp.Subscriber {
	return getQueueInstance().subscriber
}
