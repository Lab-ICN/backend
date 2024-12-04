package jwt

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Lab-ICN/backend/token-service/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
)

func Validate(token string, secret string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	}
	_token, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("parsing jwt token: %w", err)
	}
	if !_token.Valid {
		return nil, usecase.Error{Code: http.StatusUnauthorized}
	}
	return _token, nil
}
