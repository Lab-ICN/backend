package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Lab-ICN/backend/user-service/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	keyClientID = "id"
)

func BearerAuth(key string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authorization := c.Get(fiber.HeaderAuthorization)
		bearer := strings.SplitN(authorization, " ", 2)
		if bearer[0] != "Bearer" || len(bearer) != 2 {
			return &usecase.Error{
				Code:    http.StatusBadRequest,
				Message: msgInvalidBearer,
			}
		}
		token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Name}))
		if err != nil {
			return &usecase.Error{Code: http.StatusUnauthorized, Message: err.Error(), Err: err}
		}
		if !token.Valid {
			return &usecase.Error{Code: http.StatusUnauthorized}
		}
		sub, err := token.Claims.GetSubject()
		if err != nil {
			return &usecase.Error{Code: http.StatusBadRequest, Message: msgMissingSub, Err: err}
		}
		id, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			return fmt.Errorf("parsing jwt sub of %s: %w", sub, err)
		}
		c.Locals(keyClientID, id)
		return c.Next()
	}
}

func ApiKeyAuth(key string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		header := c.Get(fiber.HeaderAuthorization, "")
		if header == "" {
			return &usecase.Error{Code: http.StatusUnauthorized, Message: msgMissingAuthorization}
		}
		sent, err := base64.StdEncoding.DecodeString(header)
		if err != nil {
			return &usecase.Error{Code: http.StatusBadRequest, Err: err}
		}
		actual, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return &usecase.Error{Code: http.StatusBadRequest, Err: err}
		}
		if !bytes.Equal(sent, actual) {
			return &usecase.Error{Code: http.StatusUnauthorized, Message: msgIncorrectApiKey}
		}
		return c.Next()
	}
}
