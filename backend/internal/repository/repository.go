package repository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}
	slog.Info("database connected")
	return pool, nil
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationSQL string) error {
	_, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS _schema_version (version INT PRIMARY KEY, applied_at TIMESTAMPTZ DEFAULT NOW())`)
	if err != nil {
		return fmt.Errorf("create version table: %w", err)
	}

	var currentVersion int
	_ = pool.QueryRow(ctx, `SELECT COALESCE(MAX(version), 0) FROM _schema_version`).Scan(&currentVersion)

	if currentVersion >= 1 {
		slog.Info("migrations already applied", "version", currentVersion)
		return nil
	}

	if _, err := pool.Exec(ctx, migrationSQL); err != nil {
		return fmt.Errorf("migration 1: %w", err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO _schema_version (version) VALUES (1)`); err != nil {
		return fmt.Errorf("record migration: %w", err)
	}
	slog.Info("migration applied", "version", 1)
	return nil
}

func IsEmpty(ctx context.Context, pool *pgxpool.Pool) (bool, error) {
	var count int
	err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM pharmacies`).Scan(&count)
	if err != nil {
		return true, err
	}
	return count == 0, nil
}
