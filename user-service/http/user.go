package http

import (
	"net/http"
	"path/filepath"

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

func RegisterHandler(
	usecase usecase.IUserUsecase,
	r fiber.Router,
	validate *validator.Validate,
	logger *zap.Logger,
) {
	h := Handler{usecase, validate, logger}
	v1 := r.Group("/v1/users")
	v1.Post("/", h.Post)
	v1.Get("/:id<int>", h.Get)
	v1.Delete("/:id<int\\>", h.Delete)
}

func (h *Handler) Post(c *fiber.Ctx) error {
	if c.Get(fiber.HeaderContentType) == fiber.MIMEMultipartForm {
		filehead, err := c.FormFile(keyFile)
		if err != nil {
			return err
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
		return err
	}
	user, err := h.usecase.Register(c.Context(), payload)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(user)
}

func (h *Handler) Get(c *fiber.Ctx) error {
	_id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	id := uint64(_id)
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
