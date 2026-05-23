package comment

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type ServiceInterface interface {
	GetTrainComments(ctx context.Context, trainID uuid.UUID) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment string, train_id, user_id uuid.UUID, parentid *uuid.UUID) error
	DeleteComment(ctx context.Context, commentID uuid.UUID) error
	UpdateComment(ctx context.Context, comment string, commentID uuid.UUID) error
}

type Handler struct {
	log       *slog.Logger
	service   ServiceInterface
	validator *validator.Validate
}

func NewHandler(log *slog.Logger, service ServiceInterface) *Handler {
	valid := validator.New()
	return &Handler{
		log:       log,
		service:   service,
		validator: valid,
	}
}

func (h *Handler) GetTrainComments(w http.ResponseWriter, r *http.Request) {
	const op = "Comment.Handler.GetTrainComments"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id is required"})
		return
	}
	uuidcom, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parsing id from url")
	}
	comments, err := h.service.GetTrainComments(r.Context(), uuidcom)
	if err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.JSON(w, r, comments)
}

func (h *Handler) CreateTrainComment(w http.ResponseWriter, r *http.Request) {
	const op = "Comment.Handler.CreateTrainComment"
	log := h.log.With(slog.String("op", op))
	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}
	var req CreateCommentRequest
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
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id is required"})
		return
	}
	uuidcom, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parsing id from url")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id is required"})
		return
	}
	if err := h.service.CreateComment(r.Context(), req.Comment, uuidcom, claims.ID, req.ParentID); err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	const op = "Comment.Handler.UpdateComment"
	log := h.log.With(slog.String("op", op))

	var req UpdateCommentRequest
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
	id := chi.URLParam(r, "id")
	if id == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id is required"})
		return
	}
	uuidcom, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parsing id from url")
	}

	if err := h.service.UpdateComment(r.Context(), req.Comment, uuidcom); err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "failed to update comment"})
		return
	}
}
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	const op = "Comment.Handler.DeleteComment"
	log := h.log.With(slog.String("op", op))

	id := chi.URLParam(r, "id")
	if id == "" {
		log.Error("id is required")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id is required"})
		return
	}
	uuidcom, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed parsing id from url")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "invalid request body")
		return
	}
	if err := h.service.DeleteComment(r.Context(), uuidcom); err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.Status(r, http.StatusNoContent)
	render.JSON(w, r, "ok")
}
