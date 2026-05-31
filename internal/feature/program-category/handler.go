package programcategory

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

type ProgramCategoryServiceInterface interface {
	GetAllProgramCategory(ctx context.Context) ([]ProgramCategoryDB, error)
	CreateProgramCategory(ctx context.Context, title string) (*ProgramCategoryDB, error)
	UpdateProgramCategory(ctx context.Context, id uuid.UUID, title string) (*ProgramCategoryDB, error)
	DeleteProgramCategory(ctx context.Context, id uuid.UUID) error
}

type Handler struct {
	log       *slog.Logger
	pcservice ProgramCategoryServiceInterface
	validator *validator.Validate
}

func NewHandler(log *slog.Logger, service ProgramCategoryServiceInterface) *Handler {
	return &Handler{
		log:       log,
		pcservice: service,
		validator: validator.New(),
	}
}

// @Summary GetAllProgramCategory
// @Tags program-category
// @Description Get all program category
// @Produce json
// @Success 200 {object} []programcategory.ProgramCategoryDB "All data category programs"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /program/category [get]
func (h *Handler) GetAllProgramCategory(w http.ResponseWriter, r *http.Request) {
	const op = "ProgramCategory.Handler.GetAllProgramCategory"
	log := h.log.With(slog.String("op", op))

	categories, err := h.pcservice.GetAllProgramCategory(r.Context())
	if err != nil {
		log.Error("failed get all categories: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to get the program category")
		return
	}
	render.JSON(w, r, categories)
}

// @Summary CreateProgramCategory
// @Tags program-category
// @Description Create category program
// @Produce json
// @Param input body CreateProgramCategoryRequest true "Create body json"
// @Success 201 {object} programcategory.CreateProgramCategoryResponse "create category response"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /program/category [post]
func (h *Handler) CreateProgramCategory(w http.ResponseWriter, r *http.Request) {
	const op = "ProgramCategory.Handler.CreateProgramCategory"
	log := h.log.With(slog.String("op", op))

	var req CreateProgramCategoryRequest
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
	category, err := h.pcservice.CreateProgramCategory(r.Context(), req.Title)
	if err != nil {
		log.Error("failed create category: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to create category")
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateProgramCategoryResponse{
		Category: category,
	})
}

// @Summary UpdateProgramCategory
// @Tags program-category
// @Description Update category program by id
// @Produce json
// @Param input body UpdateProgramCategoryRequest true "Update body json"
// @Success 200 {object} programcategory.UpdateProgramCategoryResponse "update category response"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /program/category/{id} [put]
func (h *Handler) UpdateProgramCategory(w http.ResponseWriter, r *http.Request) {
	const op = "ProgramCategory.Handler.UpdateProgramCategory"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Error("failed program category id in path")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed program category id in path")
		return
	}
	parseuudi, err := uuid.Parse(id)
	if err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed program category id in path")
		return
	}
	var req UpdateProgramCategoryRequest
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
	category, err := h.pcservice.UpdateProgramCategory(r.Context(), parseuudi, req.Title)
	if err != nil {
		log.Error("failed update program category: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to update program category")
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, UpdateProgramCategoryResponse{
		Category: category,
	})
}

// @Summary DeleteProgramCategory
// @Tags program-category
// @Description delete category program by id
// @Produce json
// @Success 200 {object} programcategory.UpdateProgramCategoryResponse "delete category response"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /program/category/{id} [delete]
func (h *Handler) DeleteProgramCategory(w http.ResponseWriter, r *http.Request) {
	const op = "ProgramCategory.Handler.DeleteProgramCategory"
	log := h.log.With(slog.String("op", op))
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Error("failed program category id in path")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed program category id in path")
		return
	}
	parseuudi, err := uuid.Parse(id)
	if err != nil {
		log.Error(err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, "failed program category id in path")
		return
	}
	if err := h.pcservice.DeleteProgramCategory(r.Context(), parseuudi); err != nil {
		log.Error("failed delete program category: " + err.Error())
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, "failed to delete program category")
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, "ok")
}
