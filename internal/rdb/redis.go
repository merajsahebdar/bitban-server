package mem

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"go.giteam.ir/giteam/internal/conf"
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
