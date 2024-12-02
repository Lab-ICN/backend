package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Lab-ICN/backend/token-service/internal/usecase"
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
			return usecase.Error{
				Code:    http.StatusBadRequest,
				Message: msgInvalidBearer,
			}
		}
		token, err := jwt.Parse(bearer[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if !token.Valid {
			return usecase.Error{Code: http.StatusUnauthorized}
		}
		sub, err := token.Claims.GetSubject()
		if err != nil {
			return err
		}
		id, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			return err
		}
		c.Locals(keyClientID, id)
		return c.Next()
	}
}
