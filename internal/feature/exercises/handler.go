package exercises

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/Sanchir01/fitnow/pkg/api"
)

type ExerciseService interface {
	GetAllExercises(ctx context.Context) ([]domain.Exercise, error)
	CreateExercise(ctx context.Context, title, desc string) (*domain.Exercise, error)
}

type Handler struct {
	log       *slog.Logger
	service   ExerciseService
	validator *validator.Validate
}

func NewHandler(log *slog.Logger, service ExerciseService) *Handler {
	return &Handler{
		log:       log,
		service:   service,
		validator: validator.New(),
	}
}

type CreateExerciseRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// @Summary GetAllExercises
// @Tags exercises
// @Description Get all exercises
// @Produce json
// @Success 200 {object} []domain.Exercise "All exercises"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train/exercises [get]
func (h *Handler) GetAllExercises(w http.ResponseWriter, r *http.Request) {
	const op = "Exercises.Handler.GetAllExercises"
	log := h.log.With("op", op)

	exercises, err := h.service.GetAllExercises(r.Context())
	if err != nil {
		log.Error("failed to get exercises", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to get exercises"))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, exercises)
}

// @Summary CreateExercise
// @Tags exercises
// @Description Create new exercise
// @Accept json
// @Produce json
// @Param input body CreateExerciseRequest true "Create exercise body json"
// @Success 201 {object} domain.Exercise "Created exercise"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/exercises [post]
func (h *Handler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	const op = "Exercises.Handler.CreateExercise"
	log := h.log.With("op", op)

	claims, ok := r.Context().Value(constants.UserClaimsKey).(*constants.UserClaims)
	if !ok || claims == nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("unauthorized"))
		return
	}

	var req CreateExerciseRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode body")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		log.Error("invalid request", slog.Any("err", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}
	exercise, err := h.service.CreateExercise(r.Context(), req.Title, req.Description)
	if err != nil {
		log.Error("failed to create exercise", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to create exercise"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, exercise)
}
