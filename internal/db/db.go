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
	"go.uber.org/zap"
	"regeet.io/api/internal/conf"
	"regeet.io/api/internal/util"
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
