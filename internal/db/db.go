package db

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://localhost:5432/postgres?sslmode=disable"
	}
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func RunMigration(ctx context.Context, pool *pgxpool.Pool, sql string) error {
	// split into statements and execute one by one for clearer errors
	parts := strings.Split(sql, ";")
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" {
			continue
		}
		if _, err := pool.Exec(ctx, s); err != nil {
			return fmt.Errorf("migration failed: %w; statement: %.200s", err, s)
		}
	}
	return nil
}
