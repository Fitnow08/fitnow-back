package train

import (
	"context"
	"errors"
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
)

type TrainService interface {
	GetAllPublicTrains(ctx context.Context, param AllTrainsParams) ([]*domain.Train, error)
	GetTrainByID(ctx context.Context, id uuid.UUID) (*domain.Train, error)
	CreateTrain(ctx context.Context, req CreateTrainRequest, userID uuid.UUID) (*domain.Train, error)
	UpdateTrain(ctx context.Context, id uuid.UUID, req UpdateTrainRequest, userID uuid.UUID) (*domain.Train, error)
	DeleteTrain(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	GetUserTrains(ctx context.Context, userID uuid.UUID) ([]*domain.Train, error)
	AddUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	RemoveUserTrain(ctx context.Context, userID, trainID uuid.UUID) error
	UploadTrainImage(ctx context.Context, trainID uuid.UUID, ext, contentType string, size int64, r io.Reader) error
}

type Handler struct {
	log          *slog.Logger
	trainservice TrainService
	validator    *validator.Validate
}

func NewHandler(log *slog.Logger, trainservice TrainService) *Handler {
	return &Handler{
		log:          log,
		trainservice: trainservice,
		validator:    NewValidator(),
	}
}

// @Summary GetAllTrains
// @Tags train
// @Description Get all public trains with pagination
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size (max 100)" default(20)
// @Success 200 {object} []domain.Train "All public trains"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train [get]
func (h *Handler) GetAllTrains(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.GetAllTrains"
	log := h.log.With("op", op)
	q := r.URL.Query()
	page := 1
	if v := q.Get("page"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 {
			page = l
		}
	}
	limit := 20
	if v := q.Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil {
			limit = l
		}
		if limit > 100 {
			limit = 100
		}
	}
	allTrainsParams := AllTrainsParams{
		Page:  uint64(page),
		Limit: uint64(limit),
	}
	trains, err := h.trainservice.GetAllPublicTrains(r.Context(), allTrainsParams)
	if err != nil {
		log.Error("failed to get trains", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to get trains"))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, trains)
}

// @Summary GetTrainByID
// @Tags train
// @Description Get train by id
// @Produce json
// @Param id path string true "train id"
// @Success 200 {object} domain.Train "Train by id"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Router /train/{id} [get]
func (h *Handler) GetTrainByID(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.GetTrainByID"
	log := h.log.With("op", op)

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	train, err := h.trainservice.GetTrainByID(r.Context(), id)
	if err != nil {
		log.Error("failed to get train", slog.Any("err", err))
		if errors.Is(err, pgx.ErrNoRows) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, api.Error("train not found"))
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, api.Error("failed to get train"))
		}
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, train)
}

// @Summary CreateTrain
// @Tags train
// @Description Create new train
// @Accept json
// @Produce json
// @Param input body CreateTrainRequest true "Create train body json"
// @Success 201 {object} domain.Train "Created train"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train [post]
func (h *Handler) CreateTrain(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.CreateTrain"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("invalid token"))
		return
	}
	var req CreateTrainRequest
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
	train, err := h.trainservice.CreateTrain(r.Context(), req, claims.ID)
	if err != nil {
		log.Error("failed to create train", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to create train"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, train)
}

// @Summary UpdateTrain
// @Tags train
// @Description Update train by id
// @Accept json
// @Produce json
// @Param id path string true "train id"
// @Param input body UpdateTrainRequest true "Update train body json"
// @Success 200 {object} domain.Train "Updated train"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id} [put]
func (h *Handler) UpdateTrain(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.UpdateTrain"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	var req UpdateTrainRequest
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
	train, err := h.trainservice.UpdateTrain(r.Context(), id, req, claims.ID)
	if err != nil {
		log.Error("failed to update train", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to update train"))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, train)
}

// @Summary DeleteTrain
// @Tags train
// @Description Delete train by id
// @Produce json
// @Param id path string true "train id"
// @Success 204 "No content"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id} [delete]
func (h *Handler) DeleteTrain(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.DeleteTrain"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	if err := h.trainservice.DeleteTrain(r.Context(), id, claims.ID); err != nil {
		log.Error("failed to delete train", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to delete train"))
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// @Summary GetUserTrains
// @Tags train
// @Description Get trains of the current authorized user
// @Produce json
// @Success 200 {object} []domain.Train "User trains"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/me [get]
func (h *Handler) GetUserTrains(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.GetUserTrains"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}

	trains, err := h.trainservice.GetUserTrains(r.Context(), claims.ID)
	if err != nil {
		log.Error("failed to get user trains", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to get user trains"))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, trains)
}

// @Summary AddUserTrain
// @Tags train
// @Description Add train to the current user
// @Produce json
// @Param id path string true "train id"
// @Success 201 {object} map[string]bool "ok"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id}/add [post]
func (h *Handler) AddUserTrain(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.AddUserTrain"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}

	idStr := chi.URLParam(r, "id")
	trainID, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	if err := h.trainservice.AddUserTrain(r.Context(), claims.ID, trainID); err != nil {
		log.Error("failed to add user train", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to add train"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]bool{"ok": true})
}

// @Summary RemoveUserTrain
// @Tags train
// @Description Remove train from the current user
// @Produce json
// @Param id path string true "train id"
// @Success 204 "No content"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id}/remove [delete]
func (h *Handler) RemoveUserTrain(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.RemoveUserTrain"
	log := h.log.With("op", op)

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("user not authorized"))
		return
	}

	idStr := chi.URLParam(r, "id")
	trainID, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	if err := h.trainservice.RemoveUserTrain(r.Context(), claims.ID, trainID); err != nil {
		log.Error("failed to remove user train", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to remove train"))
		return
	}
	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// @Summary UploadTrainImage
// @Tags train
// @Description Upload train image (multipart form, field "image", max 10MB)
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "train id"
// @Param image formData file true "Train image file"
// @Success 201 {string} string "ok"
// @Failure 400 {object} domain.ErrorResponse "Bad request"
// @Failure 401 {object} domain.ErrorResponse "Unauthorized"
// @Failure 404 {object} domain.ErrorResponse "Not found"
// @Failure 500 {object} domain.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /train/{id}/image [post]
func (h *Handler) UploadTrainImage(w http.ResponseWriter, r *http.Request) {
	const op = "Train.Handler.UploadTrainImage"
	log := h.log.With("op", op)
	idStr := chi.URLParam(r, "id")
	programId, err := uuid.Parse(idStr)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid multipart form"))
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("image field required"))
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)
	contentType := header.Header.Get("Content-Type")
	if err := h.trainservice.UploadTrainImage(r.Context(), programId, ext, contentType, header.Size, file); err != nil {
		log.Error("failed to upload train", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to upload train"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}
