package rating

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type TrainService interface {
	GetAllTrainRatings(ctx context.Context) ([]RatingDB, error)
	CreateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error
	UpdateTrainRating(ctx context.Context, userid, trainid uuid.UUID, rating int) error
}
type Handler struct {
	log       *slog.Logger
	service   TrainService
	validator *validator.Validate
}

func NewHandler(log *slog.Logger, service TrainService) *Handler {
	valid := validator.New()
	return &Handler{
		log:       log,
		service:   service,
		validator: valid,
	}
}

// @Summary CreateTrainRating
// @Tags ratings
// @Description Create rating for a train by the current user
// @Accept json
// @Produce json
// @Param id path string true "train id"
// @Param input body CreateRatingRequest true "Create rating body json"
// @Success 201 {string} string "ok"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id}/ratings [post]
func (h *Handler) CreateTrainRating(w http.ResponseWriter, r *http.Request) {
	const op = "Rating.Handler.CreateTrainRating"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}
	var req CreateRatingRequest
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
		log.Error("uuid parse error", slog.Any("id", id))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "invalid request body")
		return
	}
	if err := h.service.CreateTrainRating(r.Context(), claims.ID, uuidcom, req.Rating); err != nil {
		log.Error("create train rating", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "create train rating error")
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}

// @Summary UpdateTrainRating
// @Tags ratings
// @Description Update rating for a train by the current user
// @Accept json
// @Produce json
// @Param id path string true "train id"
// @Param input body CreateRatingRequest true "Update rating body json"
// @Success 200 {string} string "ok"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id}/ratings [put]
func (h *Handler) UpdateTrainRating(w http.ResponseWriter, r *http.Request) {
	const op = "Rating.Handler.UpdateTrainRating"
	log := h.log.With("op", op)
	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}
	var req CreateRatingRequest
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
		log.Error("uuid parse error", slog.Any("id", id))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "invalid request body")
		return
	}
	if err := h.service.UpdateTrainRating(r.Context(), claims.ID, uuidcom, req.Rating); err != nil {
		log.Error("update train rating", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "update train rating error")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, "ok")
}
