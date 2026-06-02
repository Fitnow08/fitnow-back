package exercises

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/Sanchir01/fitnow/pkg/api"
)

type ExerciseService interface {
	GetAllExercises(ctx context.Context) ([]domain.Exercise, error)
	CreateExercise(ctx context.Context, title, desc string) (*domain.Exercise, error)
	SaveExerciseVideo(ctx context.Context, exercise uuid.UUID, ext, contentType string, size int64, r io.Reader) error
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

var allowedVideoTypes = map[string]struct{}{
	"video/mp4":  {},
	"video/webm": {},
	"video/avi":  {},
}

func (h *Handler) AddExerciseVideo(w http.ResponseWriter, r *http.Request) {
	const op = "Exercises.Handler.AddExerciseVideo"
	log := h.log.With("op", op)

	idStr := chi.URLParam(r, "id")
	trainID, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("failed to parse id from request", slog.Any("id", idStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 40<<20)
	if err := r.ParseMultipartForm(30 << 20); err != nil {
		log.Error("failed to parse multipart form", slog.Any("id", idStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid multipart form"))
		return
	}
	file, header, err := r.FormFile("video")
	if err != nil {
		log.Error("failed to get file from form exercises", slog.Any("id", idStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("video field required"))
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		log.Error("failed to read file head", slog.Any("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid file"))
		return
	}
	contentType := http.DetectContentType(buf[:n])
	if _, ok := allowedVideoTypes[contentType]; !ok {
		log.Error("unsupported content type", slog.String("detected", contentType))
		render.Status(r, http.StatusUnsupportedMediaType)
		render.JSON(w, r, api.Error("unsupported video format"))
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Error("failed to seek file", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("internal error"))
		return
	}

	ext := filepath.Ext(header.Filename)
	if err := h.service.SaveExerciseVideo(r.Context(), trainID, ext, contentType, header.Size, file); err != nil {
		log.Error("failed to upload program video", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to upload program video"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}
