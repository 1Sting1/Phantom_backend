package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	pkgHttp "Phantom_backend/pkg/http"
	"Phantom_backend/services/auth-service/internal/config"
	"Phantom_backend/services/auth-service/internal/models"
	"Phantom_backend/services/auth-service/internal/repository"
	"Phantom_backend/services/auth-service/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(db *sql.DB, cfg *config.Config) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	authService := services.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret)

	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgHttp.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		if err == services.ErrUserExists {
			pkgHttp.Error(w, http.StatusConflict, err.Error())
			return
		}
		pkgHttp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkgHttp.Created(w, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgHttp.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			pkgHttp.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		pkgHttp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkgHttp.Success(w, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		pkgHttp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkgHttp.Success(w, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkgHttp.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		if err == services.ErrInvalidRefreshToken {
			pkgHttp.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		pkgHttp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkgHttp.Success(w, response)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		pkgHttp.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.authService.GetUser(userID)
	if err != nil {
		pkgHttp.Error(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil {
		pkgHttp.Error(w, http.StatusNotFound, "User not found")
		return
	}

	pkgHttp.Success(w, user)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.Success(w, map[string]string{
		"status":  "ok",
		"service": "auth-service",
	})
}
