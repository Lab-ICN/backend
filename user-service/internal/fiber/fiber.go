package fiber

import (
	"errors"
	"net/http"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/Lab-ICN/backend/user-service/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func New(cfg *config.Config, log *zerolog.Logger) *fiber.App {
	return fiber.New(fiber.Config{
		DisableStartupMessage: !cfg.Development,
		ErrorHandler:          NewErrorHandler(log),
		Prefork:               true,
		RequestMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
	})
}

func NewErrorHandler(log *zerolog.Logger) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		log.Error().
			Err(err).
			Str("method", c.Method()).
			Str("endpoint", c.Path()).
			Bytes("body", c.Body()).
			Msg("error occured")
		fiberErr := new(fiber.Error)
		if errors.As(err, &fiberErr) {
			return c.SendStatus(fiberErr.Code)
		}
		uscErr := new(usecase.Error)
		if errors.As(err, &uscErr) {
			return c.Status(uscErr.Code).JSON(uscErr)
		}
		return c.SendStatus(http.StatusInternalServerError)
	}
}
