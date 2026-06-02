package user

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type UserService interface {
	GetAllUsers(ctx context.Context) ([]domain.OnlyUser, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*domain.OnlyUser, error)
}
type Handler struct {
	log     *slog.Logger
	service UserService
}

func NewHandler(log *slog.Logger, service UserService) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "Auth.Handler.GetAllUsers"
	log := h.log.With(slog.String("op", op))
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		log.Error("failed to get all users", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed get all users")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, users)
}
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	const op = "Auth.Handler.GetAllUsers"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	userid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse user id", slog.Any("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "invalid user ikd")
		return
	}
	user, err := h.service.GetUserById(r.Context(), userid)
	if err != nil {
		log.Error("failed to get user", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed get user")
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}
