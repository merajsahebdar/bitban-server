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
