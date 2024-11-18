package usecase

import (
	"context"
	"encoding/csv"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/Lab-ICN/backend/user-service/repository"
	"github.com/Lab-ICN/backend/user-service/types"
	"go.uber.org/zap"
)

const (
	colEmail uint = iota
	colUsername
	colFullname
	colIsMember
	colInternshipStartDate
)

type IUserUsecase interface {
	Register(
		ctx context.Context,
		user *types.CreateUserParams,
	) (types.User, error)
	RegisterBulkCSV(
		ctx context.Context,
		fileheader *multipart.FileHeader,
	) error
	Fetch(ctx context.Context, id uint64) (types.User, error)
	Delete(ctx context.Context, id uint64) error
}

type usecase struct {
	store  repository.IUserStorage
	logger *zap.Logger
}

func NewUserUsecase(
	store repository.IUserStorage,
	logger *zap.Logger,
) IUserUsecase {
	return &usecase{store, logger}
}

func (u *usecase) Register(
	ctx context.Context,
	user *types.CreateUserParams,
) (types.User, error) {
	id, err := u.store.Create(ctx, user)
	if err != nil {
		return types.User{}, err
	}
	return types.User{
		ID:                  id,
		Email:               user.Email,
		Username:            user.Username,
		Fullname:            user.Fullname,
		IsMember:            user.IsMember,
		InternshipStartDate: user.InternshipStartDate,
	}, nil
}

func (u *usecase) RegisterBulkCSV(
	ctx context.Context,
	fileheader *multipart.FileHeader,
) error {
	file, err := fileheader.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			u.logger.Error("failed to close file", zap.Error(err))
		}
	}()
	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return err
	}
	users := make([]types.CreateUserParams, len(rows))
	for i, cols := range rows {
		users[i].Email = cols[colEmail]
		users[i].Username = cols[colUsername]
		users[i].Fullname = cols[colFullname]
		_bool, err := strconv.ParseBool(cols[colIsMember])
		if err != nil {
			u.logger.Error("failed to parse string to boolean",
				zap.Error(err),
				zap.Strings("row", cols),
			)
			users[i].IsMember = false
		} else {
			users[i].IsMember = _bool
		}
		_time, err := time.Parse(time.DateTime, cols[colInternshipStartDate])
		if err != nil {
			u.logger.Error("failed to parse string to time.Time",
				zap.Error(err),
				zap.Strings("row", cols),
			)
			users[i].InternshipStartDate = time.Time{}
		} else {
			users[i].InternshipStartDate = _time
		}
	}
	if err := u.store.CreateBulk(ctx, users); err != nil {
		return err
	}
	return nil
}

func (u *usecase) Fetch(ctx context.Context, id uint64) (types.User, error) {
	user, err := u.store.Get(ctx, id)
	if err != nil {
		return types.User{}, err
	}
	return user.DTO(), nil
}

func (u *usecase) Delete(ctx context.Context, id uint64) error {
	return u.store.Delete(ctx, id)
}
