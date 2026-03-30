package middleware

import (
	"net/http"
	"strings"

	"Phantom_backend/pkg/jwt"
)

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	tokenService := jwt.NewTokenService(secretKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Del("X-User-ID")
			r.Header.Del("X-User-Email")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			claims, err := tokenService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Email", claims.Email)

			next.ServeHTTP(w, r)
		})
	}
}

// OptionalAuthMiddleware sets X-User-ID and X-User-Email when a valid Bearer token is present.
// Does not return 401 when token is missing or invalid (for public read + optional auth write).
func OptionalAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	tokenService := jwt.NewTokenService(secretKey)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Never trust client-supplied identity headers; only JWT may set them.
			r.Header.Del("X-User-ID")
			r.Header.Del("X-User-Email")

			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && parts[0] == "Bearer" {
					if claims, err := tokenService.ValidateToken(parts[1]); err == nil {
						r.Header.Set("X-User-ID", claims.UserID)
						r.Header.Set("X-User-Email", claims.Email)
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
