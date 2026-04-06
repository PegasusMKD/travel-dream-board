package testutil

import (
	"context"
	"testing"
	"time"

	"depgraph/internal/database"
	"depgraph/internal/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type TestDBEnv struct {
	Ctx     context.Context
	DB      *pgxpool.Pool
	Tx      pgx.Tx
	Queries *db.Queries
	Cleanup func()
}

func SetupTestDB(t *testing.T) *TestDBEnv {
	t.Helper()

	ctx := t.Context()

	// Get or create the shared container
	dsn, err := GetTestContainer(ctx)
	require.NoError(t, err)

	// Use the shared container DSN
	conn, err := database.NewPool(ctx, database.Config{
		URL:             dsn,
		MaxConns:        10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 1 * time.Minute,
	})
	require.NoError(t, err)

	// Run migrations only once (idempotent)
	err = database.RunMigrations(dsn)
	require.NoError(t, err)

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		conn.Close()
	})

	return &TestDBEnv{
		Ctx:     ctx,
		DB:      conn,
		Tx:      tx,
		Queries: db.New(tx),
		Cleanup: func() {
			_ = tx.Rollback(ctx)
			conn.Close()
		},
	}
}
