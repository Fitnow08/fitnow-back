package customMiddleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	ResponseWriter http.ResponseWriter
	status         int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw.ResponseWriter, r)

		slog.Info("time", start.Format("2006-01-02 15:04:05"))
	})
}
