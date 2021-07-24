package db

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	migrate "github.com/rubenv/sql-migrate"
	"go.giteam.ir/giteam/internal/conf"
	"go.giteam.ir/giteam/internal/util"
	"go.uber.org/zap"
)

// connectToDatabase Tries to make a connection to the database.
func connectToDatabase() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		conf.Cog.Database.Host,
		conf.Cog.Database.Port,
		conf.Cog.Database.Dbname,
		conf.Cog.Database.User,
		conf.Cog.Database.Pass,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		conf.Log.Fatal("failed to open the database", zap.String("error", err.Error()))
	}

	// Time to continue retrying
	duration := 30 * time.Second

	err = backoff.Retry(func() error {
		err = db.Ping()
		if err != nil {
			conf.Log.Warn("failed to establish a database connection, will attempt again...")
		}

		return err
	}, util.NewExponentialBackoff(duration))

	if err != nil {
		conf.Log.Fatal("failed to establish a database connection", zap.Duration("backoff", duration), zap.String("error", err.Error()))
	}

	return db
}

// dbConnectionLock
var dbConnectionLock = &sync.Mutex{}

// dbInstance Keeps a singleton instance of database connection.
var dbInstance *sql.DB

// GetDbInstance
func GetDbInstance() *sql.DB {
	if dbInstance == nil {
		dbConnectionLock.Lock()
		defer dbConnectionLock.Unlock()

		if dbInstance == nil {
			dbInstance = connectToDatabase()
		}
	}

	return dbInstance
}

// migrationLock
var migrationLock = &sync.Mutex{}

// migration Keeps the database migrations
var migration migrate.HttpFileSystemMigrationSource

// GetMigration
func GetMigration() migrate.HttpFileSystemMigrationSource {
	empty := migrate.HttpFileSystemMigrationSource{}
	if migration == empty {
		migrationLock.Lock()
		defer migrationLock.Unlock()

		if migration == empty {
			migration = migrate.HttpFileSystemMigrationSource{
				FileSystem: pkger.Dir("/migrations"),
			}
		}
	}

	return migration
}
