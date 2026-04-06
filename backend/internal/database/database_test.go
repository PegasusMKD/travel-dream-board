package database_test

import (
	"context"
	"github.com/PegasusMKD/travel-dream-board/internal/database"
	"github.com/PegasusMKD/travel-dream-board/internal/testutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	dbURL := "testUrl"
	maxConns := 25
	maxIdleConns := 5
	maxLifetime := 5 * time.Minute

	cfg := database.GetConfig(dbURL, maxConns, maxIdleConns, maxLifetime)

	assert.Equal(t, dbURL, cfg.URL)
	assert.Equal(t, maxConns, cfg.MaxConns)
	assert.Equal(t, maxIdleConns, cfg.MaxIdleConns)
	assert.Equal(t, maxLifetime, cfg.ConnMaxLifetime)
}

func TestNewPool(t *testing.T) {
	dsn, err := testutil.GetTestContainer(t.Context())
	require.NoError(t, err)

	timeoutCtx, cancel := context.WithTimeout(t.Context(), 1*time.Millisecond)
	defer cancel()

	cases := []struct {
		Name          string
		Conf          database.Config
		Ctx           context.Context
		ErrorContains string
	}{
		{
			Name:          "Improper connection string",
			Conf:          database.Config{URL: "invalid URL"},
			Ctx:           t.Context(),
			ErrorContains: "failed to parse database URL",
		},
		{
			Name:          "Improper port",
			Conf:          database.Config{URL: "postgres://user:pass@localhost:9999/nonexistent?sslmode=disable"},
			Ctx:           t.Context(),
			ErrorContains: "failed to create connection pool",
		},
		{
			Name: "Valid Case",
			Ctx:  t.Context(),
			Conf: database.Config{URL: dsn, MaxConns: 25, MaxIdleConns: 5, ConnMaxLifetime: 1 * time.Minute},
		},
		{
			Name:          "Failed ping",
			Ctx:           timeoutCtx,
			Conf:          database.Config{URL: dsn, MaxConns: 25, MaxIdleConns: 5, ConnMaxLifetime: 1 * time.Minute},
			ErrorContains: "failed to ping database",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			pool, err := database.NewPool(testCase.Ctx, testCase.Conf)
			if testCase.ErrorContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.ErrorContains)
				assert.Nil(t, pool)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, pool)
				defer pool.Close()

				assert.NoError(t, pool.Ping(testCase.Ctx))
			}
		})
	}
}

func TestPoolClose(t *testing.T) {
	dsn, err := testutil.GetTestContainer(t.Context())
	require.NoError(t, err)

	pool, err := database.NewPool(t.Context(), database.Config{URL: dsn, MaxConns: 25, MaxIdleConns: 5, ConnMaxLifetime: 1 * time.Minute})
	require.NoError(t, err)

	database.Close(pool)
}

func TestSetupNewPool(t *testing.T) {
	dsn, err := testutil.GetTestContainer(t.Context())
	require.NoError(t, err)

	cases := []struct {
		Name          string
		Conf          database.Config
		Ctx           context.Context
		ErrorContains string
	}{
		{
			Name:          "Improper connection string",
			Conf:          database.Config{URL: "invalid URL"},
			Ctx:           t.Context(),
			ErrorContains: "failed to parse database URL",
		},
		{
			Name: "Valid Case",
			Ctx:  t.Context(),
			Conf: database.Config{URL: dsn, MaxConns: 25, MaxIdleConns: 5, ConnMaxLifetime: 1 * time.Minute},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.Name, func(t *testing.T) {
			pool, err := database.SetupDatabasePool(testCase.Conf)
			if testCase.ErrorContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.ErrorContains)
				assert.Nil(t, pool)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, pool)
				defer pool.Close()

				assert.NoError(t, pool.Ping(testCase.Ctx))
			}
		})
	}
}
