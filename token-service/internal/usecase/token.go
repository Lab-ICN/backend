package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Lab-ICN/backend/token-service/internal/config"
	"github.com/Lab-ICN/backend/token-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
)

type ITokenUsecase interface {
	Generate(ctx context.Context, email string) (string, string, error)
	Refresh(ctx context.Context, id uint64) (string, error)
	Invalidate(ctx context.Context, id uint64) error
}

type usecase struct {
	store repository.ITokenStorage
	cfg   *config.Config
}

func NewTokenUsecase(store repository.ITokenStorage, cfg *config.Config) ITokenUsecase {
	return &usecase{store, cfg}
}

func (u *usecase) Generate(ctx context.Context, email string) (string, string, error) {
	id, err := u.store.GetUserID(ctx, email)
	if err != nil {
		if errors.Is(repository.ErrNoRow, err) {
			return "", "", &Error{
				Code:    http.StatusNotFound,
				Message: msgUserNotRegistered,
			}
		}
		return "", "", fmt.Errorf("fetch user id by email of %s: %w", email, err)
	}
	refresh := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Subject:  fmt.Sprint(id),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().
					UTC().
					Add(time.Duration(u.cfg.JWT.RefreshTTL) * time.Minute),
			),
		},
	)
	access := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Subject: fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(time.Now().
				UTC().
				Add(time.Duration(u.cfg.JWT.AccessTTL) * time.Minute)),
		},
	)
	refreshToken, err := refresh.SignedString([]byte(u.cfg.JWT.Key))
	if err != nil {
		return "", "", fmt.Errorf("signing refresh token: %w", err)
	}
	accessToken, err := access.SignedString([]byte(u.cfg.JWT.Key))
	if err != nil {
		return "", "", fmt.Errorf("signing access token: %w", err)
	}
	if err = u.store.CreateRefreshToken(ctx, email, refreshToken); err != nil {
		return "", "", err
	}
	return refreshToken, accessToken, nil
}

func (u *usecase) Refresh(ctx context.Context, id uint64) (string, error) {
	refreshToken, err := u.store.GetRefreshTokenByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoRow) {
			return "", Error{Code: http.StatusUnauthorized}
		}
		return "", fmt.Errorf("fetching refresh token: %w", err)
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(u.cfg.JWT.Key), err
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = u.store.DeleteRefreshToken(ctx, id)
			if err != err {
				return "", fmt.Errorf("deleting refresh token: %w", err)
			}
		}
		return "", fmt.Errorf("parsing jwt token: %w", err)
	}
	if !token.Valid {
		return "", Error{Code: http.StatusUnauthorized}
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return "", fmt.Errorf("getting jwt token subject: %w")
	}
	tokenID, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parsing jwt subject of %s to uint64: %w", tokenID, err)
	}
	if id != tokenID {
		return "", Error{Code: http.StatusUnauthorized}
	}
	access := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Subject: fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(time.Now().
				UTC().
				Add(time.Duration(u.cfg.JWT.AccessTTL) * time.Minute)),
		},
	)
	accessToken, err := access.SignedString([]byte(u.cfg.JWT.Key))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (u *usecase) Invalidate(ctx context.Context, id uint64) error {
	return u.store.DeleteRefreshToken(ctx, id)
}
