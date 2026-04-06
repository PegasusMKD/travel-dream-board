package database_test

import (
	"depgraph/internal/database"
	"depgraph/internal/testutil"
	embeds "depgraph/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMigrationsWithFS(t *testing.T) {
	dsn, err := testutil.GetTestContainer(t.Context())
	require.NoError(t, err)

	testCases := []struct {
		Name             string
		ErrorContains    string
		Directory        string
		ConnectionString string
	}{
		{Name: "Valid Case", ConnectionString: dsn, Directory: "migrations"},
		{Name: "Failing location", Directory: "test-directory", ConnectionString: dsn, ErrorContains: "failed to create iofs driver"},
		{Name: "Failing creation of source", Directory: "migrations", ConnectionString: "invalid-url", ErrorContains: "failed to create migrate instance"},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := database.RunMigrationsWithFS(
				tc.ConnectionString,
				embeds.EmbeddedMigrations,
				tc.Directory,
			)

			if tc.ErrorContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.ErrorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRunMigrations(t *testing.T) {
	dsn, err := testutil.GetTestContainer(t.Context())
	require.NoError(t, err)

	err = database.RunMigrations(dsn)
	require.NoError(t, err)
}

func TestRunMigrations_UpFailure(t *testing.T) {
	dsn, err := testutil.GetTestContainer(t.Context())
	require.NoError(t, err)

	// Create a dirty migration state to force m.Up() to fail
	pool, err := database.NewPool(t.Context(), database.Config{
		URL:             dsn,
		MaxConns:        5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 1 * time.Minute,
	})
	require.NoError(t, err)
	defer pool.Close()

	// Drop existing migration table and create a dirty one
	_, err = pool.Exec(t.Context(), `
		DROP TABLE IF EXISTS schema_migrations CASCADE;
		CREATE TABLE schema_migrations (
			version bigint PRIMARY KEY,
			dirty boolean NOT NULL
		);
		INSERT INTO schema_migrations (version, dirty) VALUES (999999, true);
	`)
	require.NoError(t, err)

	// Migration should fail due to dirty state
	err = database.RunMigrations(dsn)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to run migrations")

	// Clean-up so that we can re-run the migrations for all the other tests
	_, err = pool.Exec(t.Context(), `
		DROP TABLE IF EXISTS schema_migrations CASCADE;
	`)
	require.NoError(t, err)

}
