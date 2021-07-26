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

package queue

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"regeet.io/api/internal/conf"
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
