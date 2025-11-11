package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/0xHub/hub-api/internal/ports"
)

var (
	_ ports.HealthRepository = (*Store)(nil)
)

// Config holds pool configuration for the Postgres adapter.
type Config struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
}

// Store implements the application ports atop a Postgres database.
type Store struct {
	pool *pgxpool.Pool
}

// New creates a Store backed by a pgx connection pool.
func New(ctx context.Context, cfg Config) (*Store, error) {
	if cfg.URL == "" {
		return nil, fmt.Errorf("postgres: connection URL must be provided")
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse config: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = cfg.MaxConns
	}
	if cfg.MinConns > 0 {
		poolConfig.MinConns = cfg.MinConns
	}
	if cfg.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres: create pool: %w", err)
	}

	return &Store{pool: pool}, nil
}

// Close releases pooled connections.
func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// Ping verifies the database connection is established.
func (s *Store) Ping(ctx context.Context) error {
	if s.pool == nil {
		return fmt.Errorf("postgres: pool not initialised")
	}
	return s.pool.Ping(ctx)
}
