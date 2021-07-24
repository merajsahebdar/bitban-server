package main

import (
	"os"

	"github.com/alecthomas/kong"
	migrate "github.com/rubenv/sql-migrate"
	"go.giteam.ir/giteam/api"
	"go.giteam.ir/giteam/internal/conf"
	"go.giteam.ir/giteam/internal/controller"
	"go.giteam.ir/giteam/internal/db"
	"go.giteam.ir/giteam/internal/resolver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// MigrateUpCmd
type MigrateUpCmd struct{}

// Run Applies migrations.
func (cmd *MigrateUpCmd) Run() error {
	migrate.SetTable("migrations")
	appliedCount, err := migrate.Exec(db.GetDbInstance(), "postgres", db.GetMigration(), migrate.Up)
	if err != nil {
		conf.Log.Fatal("failed to apply migrations", zap.String("error", err.Error()))
	}

	if appliedCount > 0 {
		conf.Log.Info("migrations just applied", zap.Int("appliedCount", appliedCount))
	} else {
		conf.Log.Info("there are no migrations to apply")
	}

	return nil
}

// MigrateDownCmd
type MigrateDownCmd struct{}

// Run Drops migrations.
func (cmd *MigrateDownCmd) Run() error {
	migrate.SetTable("migrations")
	droppedCount, err := migrate.Exec(db.GetDbInstance(), "postgres", db.GetMigration(), migrate.Down)
	if err != nil {
		conf.Log.Fatal("failed to drop migrations", zap.String("error", err.Error()))
	}

	if droppedCount > 0 {
		conf.Log.Info("migrations just dropped", zap.Int("droppedCount", droppedCount))
	} else {
		conf.Log.Info("there are no migrations to drop")
	}

	return nil
}

// RunCmd
type RunCmd struct {
	Verbose bool `short:"v" default:"false" help:"Start in verbose mode."`
}

// Run Starts the app.
func (cmd *RunCmd) Run() error {
	conf.Log.Info("starting...", zap.Int("pid", os.Getpid()))

	// Provide app dependincies.
	opts := []fx.Option{
		// Queues
		// Controllers
		controller.AccountOpt,
		// Resolvers
		resolver.ConfigOpt,
		// APIs
		api.QueueOpt,
		api.HttpOpt,
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
		conf.Log.Fatal(err.Error())
	}
}
