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
	"time"

	"github.com/go-redis/cache/v8"
)

// cacheLock
var cacheLock = &sync.Mutex{}

// cacheInstance
var cacheInstance *cache.Cache

// GetCacheInstance
func GetCacheInstance() *cache.Cache {
	if cacheInstance == nil {
		cacheLock.Lock()
		defer cacheLock.Unlock()

		if cacheInstance == nil {
			cacheInstance = cache.New(&cache.Options{
				Redis:      GetDbInstance(),
				LocalCache: cache.NewTinyLFU(1000, time.Minute),
			})
		}
	}
	return cacheInstance
}
