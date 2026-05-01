package middleware

import (
	"context"
	"net/http"
	"strings"

	"pharmasense/backend/internal/services"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func Auth(authSvc *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string
			if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
				token = strings.TrimPrefix(header, "Bearer ")
			} else if q := r.URL.Query().Get("token"); q != "" {
				token = q
			}
			if token == "" {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			claims, err := authSvc.ValidateToken(token)
			if err != nil {
				http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(r *http.Request) *services.Claims {
	c, _ := r.Context().Value(ClaimsKey).(*services.Claims)
	return c
}
