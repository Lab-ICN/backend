package jwt

import (
	"errors"

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
		return nil, err
	}
	if !_token.Valid {
		return nil, errors.New("invalid token")
	}
	return _token, nil
}
