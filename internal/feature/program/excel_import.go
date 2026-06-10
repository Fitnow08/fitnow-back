package program

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sanchir01/fitnow/internal/models/domain"
	"io"
	"log/slog"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Sanchir01/fitnow/internal/feature/auth"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"github.com/Sanchir01/fitnow/pkg/api"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// ImportProgramsFromExcel читает .xlsx из r, парсит каждую строку в событие
// program.created и публикует её отдельным сообщением в Kafka.
// Возвращает количество успешно опубликованных строк.
//
// Ожидаемые заголовки первой строки (регистр и пробелы не важны):
// title, description, weeks, difficulty, category_id.
func (s *Service) ImportProgramsFromExcel(ctx context.Context, r io.Reader, userID uuid.UUID) (int, error) {
	const op = "Program.Service.ImportProgramsFromExcel"
	log := s.log.With(slog.String("op", op))

	f, err := excelize.OpenReader(r)
	if err != nil {
		log.Error("failed to open excel", slog.String("err", err.Error()))
		return 0, err
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return 0, fmt.Errorf("excel has no sheets")
	}
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Error("failed to read rows", slog.String("err", err.Error()))
		return 0, err
	}
	if len(rows) < 2 {
		return 0, fmt.Errorf("excel has no data rows")
	}

	idx := headerIndex(rows[0])
	published := 0
	for i, row := range rows[1:] {
		rowNum := i + 2 // +1 за заголовок, +1 за нумерацию с единицы

		title := cell(row, idx, "title")
		if title == "" {
			continue // пропускаем пустые строки
		}

		weeks, err := strconv.Atoi(cell(row, idx, "weeks"))
		if err != nil {
			log.Error("invalid weeks, skip row", slog.Int("row", rowNum), slog.String("err", err.Error()))
			continue
		}

		var categoryID *uuid.UUID
		if raw := cell(row, idx, "category_id"); raw != "" {
			id, err := uuid.Parse(raw)
			if err != nil {
				log.Error("invalid category_id, skip row", slog.Int("row", rowNum))
				continue
			}
			categoryID = &id
		}

		event := domain.ProgramCreatedEvent{
			Title:       title,
			Description: cell(row, idx, "description"),
			Weeks:       weeks,
			Difficulty:  domain.Level(strings.ToLower(cell(row, idx, "difficulty"))),
			CategoryID:  categoryID,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
		}

		value, err := json.Marshal(event)
		if err != nil {
			log.Error("failed to marshal row", slog.Int("row", rowNum), slog.String("err", err.Error()))
			continue
		}
		if err := s.eventservice.CreateEvent(ctx, constants.TopicProgramCreated, string(value)); err != nil {
			log.Error("failed to publish row", slog.Int("row", rowNum), slog.String("err", err.Error()))
			continue
		}
		published++
	}

	if published == 0 {
		return 0, fmt.Errorf("no valid rows published")
	}
	log.Info("excel import finished", slog.Int("published", published), slog.Int("rows", len(rows)-1))
	return published, nil
}

// headerIndex строит карту "имя колонки -> индекс" из строки заголовков.
func headerIndex(header []string) map[string]int {
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}
	return idx
}

// cell безопасно достаёт значение ячейки по имени колонки.
func cell(row []string, idx map[string]int, name string) string {
	i, ok := idx[name]
	if !ok || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

// ImportProgramsExcel — HTTP-хендлер: принимает multipart .xlsx (поле "file"),
// парсит его и публикует строки в Kafka.
func (h *Handler) ImportProgramsExcel(w http.ResponseWriter, r *http.Request) {
	const op = "Program.Handler.ImportProgramsExcel"
	log := h.log.With(slog.String("op", op))

	claims, err := auth.GetUserByHttpContext(r)
	if err != nil {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, api.Error("invalid token"))
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		log.Error("failed to parse multipart form", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid multipart form"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error("failed to get file from form", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("file field required"))
		return
	}
	defer file.Close()

	if ext := strings.ToLower(filepath.Ext(header.Filename)); ext != ".xlsx" {
		log.Error("unsupported file extension", slog.String("ext", ext))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("only .xlsx files are supported"))
		return
	}

	count, err := h.service.ImportProgramsFromExcel(r.Context(), file, claims.ID)
	if err != nil {
		log.Error("failed to import programs from excel", slog.String("err", err.Error()))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("failed to import programs"))
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, render.M{"status": "ok", "published": count})
}
