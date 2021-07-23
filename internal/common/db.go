package common

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

// dbKey
type dbKey struct{}

// ContextWithDb
func ContextWithDb(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

// GetContextDb
func GetContextDb(ctx context.Context) *sql.DB {
	return ctx.Value(dbKey{}).(*sql.DB)
}

// connectToDatabase Tries to make a connection to the database.
func connectToDatabase() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		Cog.Database.Host,
		Cog.Database.Port,
		Cog.Database.Dbname,
		Cog.Database.User,
		Cog.Database.Pass,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		Log.Fatal("failed to open the database", zap.String("error", err.Error()))
	}

	// Time to continue retrying
	duration := 30 * time.Second

	err = backoff.Retry(func() error {
		err = db.Ping()
		if err != nil {
			Log.Warn("failed to establish a database connection, will attempt again...")
		}

		return err
	}, NewExponentialBackoff(duration))

	if err != nil {
		Log.Fatal("failed to establish a database connection", zap.Duration("backoff", duration), zap.String("error", err.Error()))
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
