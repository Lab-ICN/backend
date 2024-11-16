package fiber

import (
	"net/http"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/gofiber/fiber/v2"
)

func New(cfg *config.Config) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      "acceptance-service",
		ErrorHandler: errorHandler,
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

func errorHandler(c *fiber.Ctx, err error) error {
	return nil
}
