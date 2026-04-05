package customMiddleware

import (
	"context"
	"github.com/Sanchir01/fitnow/internal/feature/auth"
	constants "github.com/Sanchir01/fitnow/internal/models/contants"
	"net/http"
	"strings"
)

func AuthMiddleware(authServiceURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authheader := r.Header.Get("Authorization")
			if authheader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			authHeaderParts := strings.Split(authheader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
				return
			}
			token := authHeaderParts[1]
			if token == "" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ParseToken(token)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), constants.UserClaimsKey, &constants.UserClaims{
				ID: claims.ID,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
