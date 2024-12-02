package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresql struct {
	conn *pgxpool.Pool
}

func NewTokenPostgreSQL(conn *pgxpool.Pool) ITokenStorage {
	return &postgresql{conn}
}

func (p *postgresql) GetUserID(ctx context.Context, email string) (uint64, error) {
	row := p.conn.QueryRow(ctx, `
        SELECT id
        FROM users
        WHERE email = $1;
    `, email)
	var id uint64
	if err := row.Scan(&id); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return 0, ErrNoRow
		}
		return 0, fmt.Errorf("scanning query result: %w", err)
	}
	return id, nil
}

func (p *postgresql) CreateRefreshToken(ctx context.Context, email string, token string) error {
	if _, err := p.conn.Exec(ctx, `
        UPDATE users
        SET refresh_token = $2
        WHERE email = $1;
    `, email, token); err != nil {
		return err
	}
	return nil
}

func (p *postgresql) DeleteRefreshToken(ctx context.Context, id uint64) error {
	if _, err := p.conn.Exec(ctx, `
        UPDATE users
        SET refresh_token = NULL
        WHERE id = $1;
    `, id); err != nil {
		return err
	}
	return nil
}

func (p *postgresql) GetRefreshTokenByID(ctx context.Context, id uint64) (string, error) {
	row := p.conn.QueryRow(ctx, `
        SELECT refresh_token
        FROM users
        WHERE id = $1;
    `, id)
	var token string
	if err := row.Scan(&token); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return "", ErrNoRow
		}
		return "", fmt.Errorf("scanning query result: %w", err)
	}
	return token, nil
}
