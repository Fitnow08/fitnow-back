package traincategory

import (
	"context"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
)

type TrainCategoryServiceInterface interface {
	GetAllTrainCategory(ctx context.Context) ([]TrainCategoryDB, error)
	CreateTrainCategory(ctx context.Context, title string) (*TrainCategoryDB, error)
	UpdateTrainCategory(ctx context.Context, id uuid.UUID, title string) (*TrainCategoryDB, error)
	DeleteTrainCategory(ctx context.Context, id uuid.UUID) error
}
type Handler struct {
	log       *slog.Logger
	tcservice TrainCategoryServiceInterface
	validator *validator.Validate
}

func NewHandler(log *slog.Logger, service TrainCategoryServiceInterface) *Handler {
	return &Handler{
		log:       log,
		tcservice: service,
		validator: validator.New(),
	}
}

// @Summary GetAllTrainCategory
// @Tags train-category
// @Description Get all train category
// @Produce json
// @Success 200 {object} []traincategory.TrainCategoryDB "All data category trains"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train/category [get]
func (h *Handler) GetAllTrainCategory(w http.ResponseWriter, r *http.Request) {
	const op = "TrainCategory.Handler.GetAllTrainCategory"
	log := h.log.With(slog.String("op", op))

	categories, err := h.tcservice.GetAllTrainCategory(r.Context())
	if err != nil {
		log.Error("failed get all categories: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to get the train category")
		return
	}
	render.JSON(w, r, categories)
}

// @Summary CreateTrainCategory
// @Tags train-category
// @Description Create category train
// @Produce json
// @Success 200 {object} traincategory.CreateTrainCategoryResponse "create category response
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train/category [post]
func (h *Handler) CreateTrainCategory(w http.ResponseWriter, r *http.Request) {
	const op = "TrainCategory.Handler.CreateTrainCategory"
	log := h.log.With(slog.String("op", op))

	var req CreateTrainCategoryRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode body")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		log.Error("invalid request", slog.Any("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}
	category, err := h.tcservice.CreateTrainCategory(r.Context(), req.Title)
	if err != nil {
		log.Error("failed create category: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to create category")
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateTrainCategoryResponse{
		Category: category,
	})
}

// @Summary UpdateTrainCategory
// @Tags train-category
// @Description Update category train by id
// @Produce json
// @Success 200 {object} traincategory.CreateTrainCategoryResponse "create category response
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train/category/{id} [put]
func (h *Handler) UpdateTrainCategory(w http.ResponseWriter, r *http.Request) {
	const op = "TrainCategory.Handler.CreateTrainCategory"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Error("failed train category id in path")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed train category id in path")
		return
	}
	parseuudi, err := uuid.Parse(id)
	if err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed train category id in path")
		return
	}
	var req CreateTrainCategoryRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode body")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}
	if err := h.validator.Struct(req); err != nil {
		log.Error("invalid request", slog.Any("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request body"))
		return
	}
	category, err := h.tcservice.UpdateTrainCategory(r.Context(), parseuudi, req.Title)
	if err != nil {
		log.Error("failed update train category: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to update train category")
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, UpdateTrainCategoryResponse{
		Category: category,
	})
}

// @Summary DeleteTrainCategory
// @Tags train-category
// @Description delete category train by id
// @Produce json
// @Success 200 {object} traincategory.CreateTrainCategoryResponse "create category response
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train/category/{id} [delete]
func (h *Handler) DeleteTrainCategory(w http.ResponseWriter, r *http.Request) {
	const op = "TrainCategory.Handler.DeleteTrainCategory"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Error("failed train category id in path")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed train category id in path")
		return
	}
	parseuudi, err := uuid.Parse(id)
	if err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed train category id in path")
		return
	}
	if err := h.tcservice.DeleteTrainCategory(r.Context(), parseuudi); err != nil {
		log.Error("failed delete train category: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to delete train category")
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, "ok")
}
