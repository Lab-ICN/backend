package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Lab-ICN/backend/token-service/internal/config"
	"github.com/Lab-ICN/backend/token-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type ITokenUsecase interface {
	Generate(ctx context.Context, email string) (string, string, error)
	Refresh(ctx context.Context, id uint64) (string, error)
	Invalidate(ctx context.Context, id uint64) error
}

type usecase struct {
	store  repository.ITokenStorage
	logger *zap.Logger
	cfg    *config.Config
}

func NewTokenUsecase(store repository.ITokenStorage, logger *zap.Logger, cfg *config.Config) ITokenUsecase {
	return &usecase{store, logger, cfg}
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

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Subject: fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(time.Now().
				UTC().
				Add(time.Duration(u.cfg.JWT.AccessTTL) * time.Minute)),
		},
	)

	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Subject: fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().
					UTC().
					Add(time.Duration(u.cfg.JWT.RefreshTTL) * time.Minute),
			),
		},
	)

	signedRefreshToken, err := refreshToken.SignedString([]byte(u.cfg.JWT.Key))
	if err != nil {
		return "", "", err
	}

	signedAccessToken, err := accessToken.SignedString([]byte(u.cfg.JWT.Key))
	if err != nil {
		return "", "", err
	}

	err = u.store.CreateRefreshToken(ctx, email, signedRefreshToken)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}

func (u *usecase) Refresh(ctx context.Context, id uint64) (string, error) {
	refreshToken, err := u.store.GetRefreshTokenByID(ctx, id)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(u.cfg.JWT.Key), err
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = u.store.DeleteRefreshToken(ctx, id)
			if err != err {
				return "", err
			}
		}
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("error token invalid")
	}

	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.RegisteredClaims{
			Subject: fmt.Sprint(id),
			ExpiresAt: jwt.NewNumericDate(time.Now().
				UTC().
				Add(time.Duration(u.cfg.JWT.AccessTTL) * time.Minute)),
		},
	)

	signedAccessToken, err := accessToken.SignedString([]byte(u.cfg.JWT.Key))
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func (u *usecase) Invalidate(ctx context.Context, id uint64) error {
	return u.store.DeleteRefreshToken(ctx, id)
}
