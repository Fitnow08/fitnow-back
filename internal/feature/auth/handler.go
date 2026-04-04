package auth

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*domain.User, error)
}
type Handler struct {
	log         *slog.Logger
	autgservice AuthService
	validator   *validator.Validate
}

func NewHandler(log *slog.Logger, autgservice AuthService) *Handler {
	valid := NewValidator()
	return &Handler{
		log:         log,
		autgservice: autgservice,
		validator:   valid,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "Auth.Handler.Login"
	log := h.log.With(slog.String("op", op))

	log.Info("login")
	render.JSON(w, r, "ok")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "Auth.Handler.Register"
	log := h.log.With(slog.String("op", op))
	log.Info("register")
	var req RegisterRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode body register")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "invalid request body")
		return
	}
	if err := h.validator.Struct(req); err != nil {
		log.Error("invalid request", slog.Any("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}

	user, err := h.autgservice.Register(r.Context(), req)
	if err != nil {
		log.Error("failed to register user", slog.Any("err", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed register")
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}
