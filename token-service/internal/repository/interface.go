package repository

import "context"

type ITokenStorage interface {
	GetUserID(ctx context.Context, email string) (uint64, error)
	CreateRefreshToken(ctx context.Context, email, token string) error
	DeleteRefreshToken(ctx context.Context, id uint64) error
	GetRefreshTokenByID(ctx context.Context, id uint64) (string, error)
}
