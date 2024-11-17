package repository

import (
	"context"

	"github.com/Lab-ICN/backend/user-service/types"
)

type IUserStorage interface {
	Create(ctx context.Context, user *types.CreateUserParams) (int64, error)
	CreateBulk(ctx context.Context, users []types.CreateUserParams) error
	List(ctx context.Context) ([]User, error)
	ListPassed(ctx context.Context, year uint) ([]User, error)
	Get(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Delete(ctx context.Context, id int64) error
}
