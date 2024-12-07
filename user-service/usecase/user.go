package usecase

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/Lab-ICN/backend/user-service/repository"
	"github.com/Lab-ICN/backend/user-service/types"
	"github.com/rs/zerolog"
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
	store repository.IUserStorage
	log   *zerolog.Logger
}

func NewUserUsecase(
	store repository.IUserStorage,
	log *zerolog.Logger,
) IUserUsecase {
	return &usecase{store, log}
}

func (u *usecase) Register(
	ctx context.Context,
	user *types.CreateUserParams,
) (types.User, error) {
	id, err := u.store.Create(ctx, user)
	if err != nil {
		if errors.Is(repository.ErrDuplicateRow, err) {
			return types.User{}, &Error{
				Code:    http.StatusConflict,
				Message: msgUserExist,
			}
		}
		return types.User{}, fmt.Errorf("register user: %w", err)
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
		return fmt.Errorf("open csv file header: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			u.log.Error().Err(err).Msg("closing file buffer")
		}
	}()
	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("read csv file content: %w", err)
	}
	users := make([]types.CreateUserParams, len(rows))
	for i, cols := range rows {
		users[i].Email = cols[colEmail]
		users[i].Username = cols[colUsername]
		users[i].Fullname = cols[colFullname]
		_bool, err := strconv.ParseBool(cols[colIsMember])
		if err != nil {
			return fmt.Errorf("parse value %s to boolean: %w", cols[colIsMember], err)
		}
		users[i].IsMember = _bool
		_time, err := time.Parse(time.RFC3339, cols[colInternshipStartDate])
		if err != nil {
			return fmt.Errorf(
				"parse value %s to RFC3339: %w",
				cols[colInternshipStartDate],
				err,
			)
		}
		users[i].InternshipStartDate = _time
	}
	if err := u.store.CreateBulk(ctx, users); err != nil {
		if errors.Is(repository.ErrDuplicateRow, err) {
			return &Error{
				Code:    http.StatusConflict,
				Message: msgUserExist,
			}
		}
		return fmt.Errorf("register user: %w", err)
	}
	return nil
}

func (u *usecase) Fetch(ctx context.Context, id uint64) (types.User, error) {
	user, err := u.store.Get(ctx, id)
	if err != nil {
		if errors.Is(repository.ErrNoRow, err) {
			return types.User{}, &Error{
				Code:    http.StatusNotFound,
				Message: msgUserNotFound,
			}
		}
		return types.User{}, fmt.Errorf("fetch user by id: %w", err)
	}
	return user.DTO(), nil
}

func (u *usecase) Delete(ctx context.Context, id uint64) error {
	return u.store.Delete(ctx, id)
}
