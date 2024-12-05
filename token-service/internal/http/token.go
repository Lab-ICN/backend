package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Lab-ICN/backend/token-service/internal/config"
	_jwt "github.com/Lab-ICN/backend/token-service/internal/jwt"
	"github.com/Lab-ICN/backend/token-service/internal/usecase"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/idtoken"
)

type Handler struct {
	usecase  usecase.ITokenUsecase
	cfg      *config.Config
	validate *validator.Validate
}

func RegisterHandlers(
	usecase usecase.ITokenUsecase,
	cfg *config.Config,
	r fiber.Router,
	validate *validator.Validate,
) {
	h := Handler{usecase, cfg, validate}
	v1 := r.Group("/v1/tokens")
	v1.Post("/", h.GenerateHandler)
	// FIXME: method patch makes panic
	v1.Put("/self", h.RefreshHandler)
	v1.Delete("/self", BearerAuth(cfg.JWT.Key), h.InvalidateHandler)
}

func (h *Handler) GenerateHandler(c *fiber.Ctx) error {
	payload := new(struct {
		Token string `json:"token"`
	})
	if err := c.BodyParser(payload); err != nil {
		return &usecase.Error{Code: fiber.StatusBadRequest}
	}
	claims, err := idtoken.Validate(c.Context(), payload.Token, h.cfg.GoogleClientID)
	if err != nil {
		return &usecase.Error{Code: http.StatusUnauthorized, Err: err}
	}
	refresh, access, err := h.usecase.Generate(c.Context(), claims.Claims["email"].(string))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"refreshToken": refresh,
		"accessToken":  access,
	})
}

func (h *Handler) RefreshHandler(c *fiber.Ctx) error {
	payload := new(struct {
		Token string `json:"refreshToken"`
	})
	if err := c.BodyParser(payload); err != nil {
		return &usecase.Error{Code: fiber.StatusBadRequest}
	}
	token, err := _jwt.Validate(payload.Token, h.cfg.JWT.Key)
	if err != nil {
		return err
	}
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return fmt.Errorf("getting jwt token subject: %w")
	}
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing jwt subject of %s to uint64: %w", id, err)
	}
	access, err := h.usecase.Refresh(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"accessToken": access})
}

func (h *Handler) InvalidateHandler(c *fiber.Ctx) error {
	id, ok := c.Locals(keyClientID).(uint64)
	if !ok {
		return &usecase.Error{Code: http.StatusInternalServerError}
	}
	if err := h.usecase.Invalidate(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(http.StatusOK)
}
