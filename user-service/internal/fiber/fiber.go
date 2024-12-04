package fiber

import (
	"errors"
	"net/http"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/Lab-ICN/backend/user-service/usecase"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.Logger) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      "acceptance-service",
		ErrorHandler: NewErrorHandler(logger),
		Prefork:      true,
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

func NewErrorHandler(logger *zap.Logger) func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		logger.Error("error occured", zap.Error(err))
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
