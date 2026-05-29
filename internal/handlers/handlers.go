package handlers

import (
	"github.com/Sanchir01/fitnow/internal/app"
	customMiddleware "github.com/Sanchir01/fitnow/internal/handlers/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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
			customMiddleware.RecoverMiddleware,
			middleware.RequestID,
			limiter.RateLimitMiddleware,
		)
		r.Get("/ws", handlers.ChatHandler.WsHandler)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", handlers.AuthHandler.Register)
			r.Post("/login", handlers.AuthHandler.Login)
			r.Post("/token", handlers.AuthHandler.NewTokens)
			r.Post("/verify", handlers.AuthHandler.VerifyAccount)
			r.Post("/verify-resend", handlers.AuthHandler.ResendVerifyCode)
			r.Post("/password/reset", handlers.AuthHandler.ResetPassword)
			r.Post("/password/reset/confirm", handlers.AuthHandler.ConfirmResetPassword)
		})
		r.Route("/train", func(r chi.Router) {
			r.Get("/", handlers.TrainHandler.GetAllTrains)
			r.Get("/exercises", handlers.ExercisesHandler.GetAllExercises)
			r.Get("/{id}", handlers.TrainHandler.GetTrainByID)

			r.Group(func(r chi.Router) {
				r.Use(customMiddleware.AuthMiddleware(""))
				r.Route("/{id}", func(r chi.Router) {
					r.Post("/image", handlers.TrainHandler.UploadTrainImage)

					r.Route("/ratings", func(r chi.Router) {
						r.Post("/", handlers.RatingHandler.CreateTrainRating)
						r.Put("/", handlers.RatingHandler.UpdateTrainRating)
					})
				})
				r.Post("/", handlers.TrainHandler.CreateTrain)

				r.Post("/exercises", handlers.ExercisesHandler.CreateExercise)
				r.Get("/me", handlers.TrainHandler.GetUserTrains)

				r.Put("/{id}", handlers.TrainHandler.UpdateTrain)
				r.Delete("/{id}", handlers.TrainHandler.DeleteTrain)
				r.Post("/{id}/add", handlers.TrainHandler.AddUserTrain)
				r.Delete("/{id}/remove", handlers.TrainHandler.RemoveUserTrain)
			})
			r.Route("/comments", func(r chi.Router) {
				r.Get("/{train-id}", handlers.CommentHandler.GetTrainComments)
				r.Group(func(r chi.Router) {
					r.Use(customMiddleware.AuthMiddleware(""))
					r.Delete("/{id}", handlers.CommentHandler.DeleteComment)
					r.Put("/{id}", handlers.CommentHandler.UpdateComment)
					r.Post("/{train-id}", handlers.CommentHandler.CreateTrainComment)
				})

			})
			r.Route("/category", func(r chi.Router) {
				r.Get("/", handlers.TrainCategoryHandler.GetAllTrainCategory)
				r.Post("/", handlers.TrainCategoryHandler.CreateTrainCategory)
				r.Put("/{id}", handlers.TrainCategoryHandler.UpdateTrainCategory)
				r.Delete("/{id}", handlers.TrainCategoryHandler.DeleteTrainCategory)
			})
		})

	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
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
