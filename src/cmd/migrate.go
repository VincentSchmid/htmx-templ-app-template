//go:build dev

package cmd

import (
	"database/sql"
	"fmt"

	"github.com/VincentSchmid/htmx-templ-app-template/pkg/appconfig"
	database "github.com/VincentSchmid/htmx-templ-app-template/pkg/db"
	"github.com/VincentSchmid/htmx-templ-app-template/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func initializeMigrations(db *sql.DB, config appconfig.Database) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		config.GetMigrationDirectory(),
		config.GetDriver(),
		driver,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migration, nil
}

var (
	migration *migrate.Migrate
	version   int

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Run database migrations`,
		Args:  cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			config := appconfig.Config.Database
			logger.Log.Info("Initializing migrations for database", zap.String("database", config.GetConnectionString()))
			dbProvider := database.NewDbProvider(config, false)
			migration, err = initializeMigrations(dbProvider.Db, config)
			if err != nil {
				return fmt.Errorf("failed to create migration instance: %w", err)
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Apply all available migrations",
		Long:  `Apply all available migrations`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Log.Info("Applying migrations")
			if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
				return err
			}

			return nil
		},
	}

	downCmd = &cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		Long:  `Rollback the last migration`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Log.Info("Rolling back migrations")
			if err := migration.Down(); err != nil && err != migrate.ErrNoChange {
				return err
			}

			return nil
		},
	}

	dropCmd = &cobra.Command{
		Use:   "drop",
		Short: "Drop all tables",
		Long:  `Drop all tables`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Log.Info("Dropping all tables")
			if err := migration.Drop(); err != nil && err != migrate.ErrNoChange {
				return err
			}

			return nil
		},
	}

	forceCmd = &cobra.Command{
		Use:   "force",
		Short: "Force a specific version of migration",
		Long:  `Force a specific version of migration`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := migration.Force(version); err != nil && err != migrate.ErrNoChange {
				return err
			}

			return nil
		},
	}
)

func init() {
	forceCmd.Flags().IntVarP(&version, "version", "v", 0, "Version to force migration to")
	migrateCmd.AddCommand(upCmd)
	migrateCmd.AddCommand(downCmd)
	migrateCmd.AddCommand(dropCmd)
	migrateCmd.AddCommand(forceCmd)
	rootCmd.AddCommand(migrateCmd)
}
