package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Lab-ICN/backend/user-service/internal/config"
	"github.com/Lab-ICN/backend/user-service/types"
	"github.com/Lab-ICN/backend/user-service/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

const (
	keyFile = "attachment"
)

type Handler struct {
	usecase  usecase.IUserUsecase
	validate *validator.Validate
	logger   *zap.Logger
}

func RegisterHandlers(
	usecase usecase.IUserUsecase,
	cfg *config.Config,
	r fiber.Router,
	validate *validator.Validate,
	logger *zap.Logger,
) {
	h := Handler{usecase, validate, logger}
	v1 := r.Group("/v1/users")
	v1.Get("/self", BearerAuth(cfg.JwtKey), h.Get)
	v1.Post("/", ApiKeyAuth(cfg.ApiKey), h.Post)
	v1.Delete("/:id<int>", ApiKeyAuth(cfg.ApiKey), h.Delete)
}

func (h *Handler) Post(c *fiber.Ctx) error {
	if strings.Contains(c.Get(fiber.HeaderContentType), fiber.MIMEMultipartForm) {
		filehead, err := c.FormFile(keyFile)
		if err != nil {
			return &usecase.Error{
				Code:    http.StatusBadRequest,
				Message: msgMissingAttachment,
				Err:     err,
			}
		}
		if filepath.Ext(filehead.Filename) != ".csv" {
			return &usecase.Error{
				Code:    http.StatusUnprocessableEntity,
				Message: msgMustCSV,
			}
		}
		if err := h.usecase.RegisterBulkCSV(c.Context(), filehead); err != nil {
			return err
		}
		return c.SendStatus(http.StatusCreated)
	}
	payload := new(types.CreateUserParams)
	if err := c.BodyParser(payload); err != nil {
		return &usecase.Error{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}
	user, err := h.usecase.Register(c.Context(), payload)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(user)
}

func (h *Handler) Get(c *fiber.Ctx) error {
	id, ok := c.Locals(keyClientID).(uint64)
	if !ok {
		return &usecase.Error{Code: http.StatusInternalServerError}
	}
	user, err := h.usecase.Fetch(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(user)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	_id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	id := uint64(_id)
	if err := h.usecase.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
