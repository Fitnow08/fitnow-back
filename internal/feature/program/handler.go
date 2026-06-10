package program

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
)

type ProgramService interface {
	CreateProgram(ctx context.Context, title string, description string, weeks int, level domain.Level, categoryID *uuid.UUID, user_id uuid.UUID) (*domain.Program, error)
	GetAllProgramAndTrainsCount(ctx context.Context, categoryId uuid.UUID, search string) ([]domain.ProgramAndTrainsCount, error)
	UploadProgramImage(ctx context.Context, programID uuid.UUID, ext, contentType string, size int64, r io.Reader) error
	UploadAllProgramTrains(ctx context.Context, programId uuid.UUID, trains []ProgramTrainInput) error
	GetProgramTrains(ctx context.Context, programID uuid.UUID) (*domain.ProgramAndTrain, error)
	ImportProgramsFromExcel(ctx context.Context, r io.Reader, userID uuid.UUID) (int, error)
	PublishProgramCreated(ctx context.Context, userID uuid.UUID, req CreateProgramRequest) error
}
type Handler struct {
	log       *slog.Logger
	service   ProgramService
	validator *validator.Validate
}

func NewHandler(log *slog.Logger, service ProgramService) *Handler {
	valid := auth.NewValidator()
	return &Handler{
		log:       log,
		service:   service,
		validator: valid,
	}
}
func (h *Handler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.CreateProgram"
	log := h.log.With(slog.String("op", op))
	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("invalid token"))
		return
	}
	var req CreateProgramRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode body verify account")
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

	program, err := h.service.CreateProgram(r.Context(), req.Title, req.Description, req.Weeks, req.Difficulty, req.CategoryID, claims.ID)
	if err != nil {
		log.Error("failed to create program", slog.Any("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, render.M{"status": "failed create program"})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, program)
}

func (h *Handler) GetAllPrograms(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.GetAllPrograms"
	log := h.log.With(slog.String("op", op))

	query := r.URL.Query()
	categoryId := query.Get("category_id")

	categoryuuid, err := uuid.Parse(categoryId)
	if err != nil {
		categoryuuid = uuid.Nil
	}
	search := query.Get("search")
	programs, err := h.service.GetAllProgramAndTrainsCount(r.Context(), categoryuuid, search)
	if err != nil {
		log.Error("failed to get all programs", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"status": "failed get all programs"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, programs)
}

func (h *Handler) AddProgramImage(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.AddProgramImage"
	log := h.log.With(slog.String("op", op))
	idStr := chi.URLParam(r, "id")
	trainID, err := uuid.Parse(idStr)
	if err != nil {
		log.Error("failed to parse id from request", slog.Any("id", idStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		log.Error("failed to parse multipart form", slog.Any("id", idStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid multipart form"))
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Error("failed to get file from form", slog.Any("id", idStr), "err", err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("image field required"))
		return
	}
	defer file.Close()
	ext := filepath.Ext(header.Filename)
	contentType := header.Header.Get("Content-Type")
	if err := h.service.UploadProgramImage(r.Context(), trainID, ext, contentType, header.Size, file); err != nil {
		log.Error("failed to upload program image", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to upload program image"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}

func (h *Handler) UploadAllProgramTrains(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.UploadAllProgramTrains"
	log := h.log.With(slog.String("op", op))

	id := chi.URLParam(r, "id")
	trainID, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id from request", slog.Any("id", id))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}

	var req AddProgramTrainsRequest
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

	if err := h.service.UploadAllProgramTrains(r.Context(), trainID, req.Trains); err != nil {
		log.Error("failed to upload program trains", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"status": "failed upload program trains"})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}

func (h *Handler) GetProgramsTrains(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.GetProgramsTrains"
	log := h.log.With(slog.String("op", op))

	id := chi.URLParam(r, "id")
	programuuid, err := uuid.Parse(id)
	if err != nil {
		log.Error("failed to parse id from request", slog.Any("id", id), "err", err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid id"))
		return
	}
	data, err := h.service.GetProgramTrains(r.Context(), programuuid)
	if err != nil {
		log.Error("failed to get program trains", slog.Any("err", err.Error()))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"status": "failed get program trains"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)

}

func (h *Handler) CreateProgramToKafka(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.CreateProgramToKafka"
	log := h.log.With(slog.String("op", op))

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("invalid token"))
		return
	}
	var req CreateProgramRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode body send message to kafka", slog.Any("err", err.Error()))
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
	if err := h.service.PublishProgramCreated(r.Context(), claims.ID, req); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, render.M{"status": "failed create program"})
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, "ok")
}
