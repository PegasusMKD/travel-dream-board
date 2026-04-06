package database

import (
	"fmt"
	"io/fs"
	embeds "depgraph/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(databaseURL string) error {
	return RunMigrationsWithFS(databaseURL, embeds.EmbeddedMigrations, "migrations")
}

func RunMigrationsWithFS(databaseURL string, migrationFS fs.FS, dir string) error {
	d, err := iofs.New(migrationFS, dir)
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
