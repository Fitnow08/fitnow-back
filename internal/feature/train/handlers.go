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
	"log/slog"
	"net/http"
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
