package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Lab-ICN/backend/user-service/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxRecords = 500

type postgresql struct {
	conn *pgxpool.Pool
}

func NewUserPostgreSQL(conn *pgxpool.Pool) IUserStorage {
	return &postgresql{conn}
}

func (p *postgresql) Create(ctx context.Context, user *types.CreateUserParams) (uint64, error) {
	var id uint64
	if err := p.conn.QueryRow(ctx, `
        INSERT INTO users ("email", "username", "fullname", "is_member", "internship_start_date")
        VALUES (@email, @username, @fullname, @is_member, @internship_start_date)
        RETURNING id`,
		pgx.NamedArgs{
			"email":                 user.Email,
			"username":              user.Username,
			"fullname":              user.Fullname,
			"is_member":             user.IsMember,
			"internship_start_date": user.InternshipStartDate,
		},
	).Scan(&id); err != nil {
		pgErr := new(pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrDuplicateRow
		}
		return 0, fmt.Errorf("inserting user for email %s: %w", user.Email, err)
	}
	return id, nil
}

func (p *postgresql) CreateBulk(ctx context.Context, users []types.CreateUserParams) error {
	affected, err := p.conn.CopyFrom(
		ctx,
		pgx.Identifier{"users"},
		[]string{"email", "username", "fullname", "is_member", "internship_start_date"},
		pgx.CopyFromSlice(len(users), func(i int) ([]interface{}, error) {
			return []interface{}{
				users[i].Email,
				users[i].Username,
				users[i].Fullname,
				users[i].IsMember,
				users[i].InternshipStartDate,
			}, nil
		}),
	)
	if err != nil {
		pgErr := new(pgconn.PgError)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateRow
		}
		return fmt.Errorf("inserting users: %w", err)
	}
	if affected == 0 {
		return fmt.Errorf("creating bulk user: %w", ErrNoRowAffected)
	}
	return nil
}

func (p *postgresql) List(ctx context.Context) ([]User, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT
			id,
			email,
			username,
			fullname,
			is_member,
			internship_start_date
		FROM users
		ORDER BY created_at
		LIMIT $1`, maxRecords,
	)
	if err != nil {
		return nil, fmt.Errorf("selecting users: %w", err)
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, fmt.Errorf("parsing users: %w", err)
	}
	return users, nil
}

func (p *postgresql) ListPassed(ctx context.Context, year uint) ([]User, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT 
			id,
			email,
			username,
			fullname,
			is_member,
			internship_start_date
		FROM users
		WHERE is_member = TRUE AND
		EXTRACT(YEAR FROM internship_start_date) = $1
		ORDER BY created_at
		LIMIT $2`, year, maxRecords,
	)
	if err != nil {
		return nil, fmt.Errorf("selecting members in year %d: %w", year, err)
	}
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, fmt.Errorf("parsing users: %w", err)
	}
	return users, nil
}

func (p *postgresql) Get(ctx context.Context, id uint64) (User, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT 
			id,
			email,
			username,
			fullname,
			is_member,
			internship_start_date
		FROM users WHERE id = $1`, id)
	if err != nil {
		return User{}, fmt.Errorf("selecting user for id %d: %w", id, err)
	}
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return User{}, ErrNoRow
		}
		return User{}, fmt.Errorf("parsing user: %w", err)
	}
	return user, nil
}

func (p *postgresql) GetByEmail(ctx context.Context, email string) (User, error) {
	rows, err := p.conn.Query(ctx, `
		SELECT 
			id,
			email,
			username,
			fullname,
			is_member,
			internship_start_date
		FROM users WHERE email = $1`, email)
	if err != nil {
		return User{}, fmt.Errorf("selecting user for email %s: %w", email, err)
	}
	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return User{}, ErrNoRow
		}
		return User{}, fmt.Errorf("parsing user: %w", err)
	}
	return user, nil
}

func (p *postgresql) Delete(ctx context.Context, id uint64) error {
	_, err := p.conn.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("deleting user for id %d: %w", id, err)
	}
	return nil
}
