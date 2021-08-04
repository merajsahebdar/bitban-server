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

package orm

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	_ "github.com/lib/pq"
	"github.com/markbates/pkger"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/zap"
	"regeet.io/api/internal/cfg"
	"regeet.io/api/internal/pkg/util"
)

// connectToDatabase Tries to make a connection to the database.
func connectToDatabase() *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Cog.Database.Host,
		cfg.Cog.Database.Port,
		cfg.Cog.Database.Dbname,
		cfg.Cog.Database.User,
		cfg.Cog.Database.Pass,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		cfg.Log.Fatal("failed to open the database", zap.String("error", err.Error()))
	}

	if cfg.CurrentEnv == cfg.Prod {
		// Time to continue retrying
		duration := 30 * time.Second

		err = backoff.Retry(func() error {
			err = db.Ping()
			if err != nil {
				cfg.Log.Warn("failed to establish a database connection, will attempt again...")
			}

			return err
		}, util.NewExponentialBackoff(duration))

		if err != nil {
			cfg.Log.Fatal("failed to establish a database connection", zap.Duration("backoff", duration), zap.String("error", err.Error()))
		}
	}

	return db
}

// dbConnectionMutex
var dbConnectionMutex = &sync.Mutex{}

// dbInstance Keeps a singleton instance of database connection.
var dbInstance *sql.DB

// GetDbInstance
func GetDbInstance() *sql.DB {
	if dbInstance == nil {
		dbConnectionMutex.Lock()
		defer dbConnectionMutex.Unlock()

		if dbInstance == nil {
			dbInstance = connectToDatabase()
		}
	}

	return dbInstance
}

// bunInstanceMutex
var bunInstanceMutex = &sync.Mutex{}

// bunInstance
var bunInstance *bun.DB

// GetBunInstance
func GetBunInstance() *bun.DB {
	if bunInstance == nil {
		bunInstanceMutex.Lock()
		defer bunInstanceMutex.Unlock()

		if bunInstance == nil {
			bunInstance = bun.NewDB(
				GetDbInstance(),
				pgdialect.New(),
			)
		}
	}

	return bunInstance
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
				FileSystem: pkger.Dir("/internal/pkg/orm/migration"),
			}
		}
	}

	return migration
}

// MigrateUp
func MigrateUp() (int, error) {
	migrate.SetTable("migrations")
	return migrate.Exec(GetDbInstance(), "postgres", GetMigration(), migrate.Up)
}

// MigrateDown
func MigrateDown(max int) (int, error) {
	migrate.SetTable("migrations")
	return migrate.ExecMax(GetDbInstance(), "postgres", GetMigration(), migrate.Down, max)
}
