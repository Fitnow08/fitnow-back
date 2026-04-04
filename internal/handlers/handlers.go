package handlers

import (
	"github.com/Sanchir01/fitnow/internal/app"
	customMiddleware "github.com/Sanchir01/fitnow/internal/handlers/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"os"
)

func StartHttpHandlers(handlers *app.Handlers) http.Handler {
	r := chi.NewRouter()
	limiter := customMiddleware.NewIPRateLimiter(5, 10)
	StartCors(r)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(
			limiter.RateLimitMiddleware,
			middleware.RequestID,
			customMiddleware.RecoverMiddleware,
		)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.AuthHandler.Register)
		})
	})
	return r
}

func StartCors(r *chi.Mux) {
	allowedOrigin := os.Getenv("FRONTEND_URL")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3020"
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}
