// Package postgres contains reusable PostgreSQL driver logic
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize = 1
	_defaultConnTimeout = time.Second
)

type Postgres struct {
	maxPoolSize int
	connTimeout time.Duration
	Pool        *pgxpool.Pool
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize: _defaultMaxPoolSize,
		connTimeout: _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres.pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	// Try to connect immediately to validate connection
	pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres.NewPostgres.failed to connect: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), pg.connTimeout)
	defer cancel()

	if err = pg.Pool.Ping(ctx); err != nil {
		pg.Pool.Close()
		return nil, fmt.Errorf("postgres.NewPostgres.failed to ping: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
