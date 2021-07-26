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

package mem

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"regeet.io/api/internal/conf"
)

// dbLock
var dbLock = &sync.Mutex{}

// dbInstance
var dbInstance *redis.Client

// GetDbInstance
func GetDbInstance() *redis.Client {
	if dbInstance == nil {
		dbLock.Lock()
		defer dbLock.Unlock()

		if dbInstance == nil {
			if opt, err := redis.ParseURL(conf.Cog.Redis.Url); err != nil {
				conf.Log.Fatal(err.Error())
			} else {
				dbInstance = redis.NewClient(opt)
			}
		}
	}

	return dbInstance
}
