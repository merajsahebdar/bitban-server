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

package main

import (
	"os"

	"github.com/alecthomas/kong"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"regeet.io/api/internal/app/api"
	"regeet.io/api/internal/app/controller"
	"regeet.io/api/internal/app/resolver"
	"regeet.io/api/internal/cfg"
	"regeet.io/api/internal/pkg/orm"
)

// MigrateUpCmd
type MigrateUpCmd struct{}

// Run Applies migrations.
func (cmd *MigrateUpCmd) Run() error {
	if appliedCount, err := orm.MigrateUp(); err != nil {
		cfg.Log.Fatal("failed to apply migrations", zap.String("error", err.Error()))
	} else if appliedCount > 0 {
		cfg.Log.Info("migrations just applied", zap.Int("appliedCount", appliedCount))
	} else {
		cfg.Log.Info("there are no migrations to apply")
	}

	return nil
}

// MigrateDownCmd
type MigrateDownCmd struct{}

// Run Drops the latest migration.
func (cmd *MigrateDownCmd) Run() error {
	if droppedCount, err := orm.MigrateDown(1); err != nil {
		cfg.Log.Fatal("failed to drop migrations", zap.String("error", err.Error()))
	} else if droppedCount > 0 {
		cfg.Log.Info("migrations just dropped", zap.Int("droppedCount", droppedCount))
	} else {
		cfg.Log.Info("there are no migrations to drop")
	}

	return nil
}

// RunCmd
type RunCmd struct {
	Verbose bool `short:"v" default:"false" help:"Start in verbose mode."`
}

// Run Starts the app.
func (cmd *RunCmd) Run() error {
	cfg.Log.Info("starting...", zap.Int("pid", os.Getpid()))

	// Provide app dependincies.
	opts := []fx.Option{
		// Controllers
		controller.AccountOpt,
		controller.RepoOpt,
		// Resolvers
		resolver.ConfigOpt,
		// APIs
		api.EchoOpt,
		api.SshOpt,
	}

	// Provide fx.NopLogger if it is not running in verbose mode.
	if !cmd.Verbose {
		opts = append(opts, fx.NopLogger)
	}

	// Run!
	fx.New(
		fx.Options(opts...),
	).Run()

	return nil
}

// CLI
var CLI struct {
	Migrate struct {
		Up   MigrateUpCmd   `cmd:"up" help:"Apply all migrations."`
		Down MigrateDownCmd `cmd:"down" help:"Drop migrations."`
	} `cmd:"migrate" help:"Run the migrator."`
	Run RunCmd `cmd:"run" help:"Run the app."`
}

// main
func main() {
	if err := kong.Parse(&CLI).Run(&CLI); err != nil {
		cfg.Log.Fatal(err.Error())
	}
}
