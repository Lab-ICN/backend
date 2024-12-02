package postgresql

import (
	"context"
	"fmt"

	"github.com/Lab-ICN/backend/token-service/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConn(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	uri := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Address,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)
	_cfg, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}
	conn, err := pgx.ConnectConfig(context.Background(), _cfg)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	uri := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgreSQL.Username,
		cfg.PostgreSQL.Password,
		cfg.PostgreSQL.Address,
		cfg.PostgreSQL.Port,
		cfg.PostgreSQL.Database,
	)
	_cfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, _cfg)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
