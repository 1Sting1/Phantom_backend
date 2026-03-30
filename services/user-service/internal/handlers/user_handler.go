package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	pkgHttp "Phantom_backend/pkg/http"
	"Phantom_backend/services/user-service/internal/config"
	"Phantom_backend/services/user-service/internal/models"
	"Phantom_backend/services/user-service/internal/repository"
)

type UserHandler struct {
	profileRepo  *repository.ProfileRepository
	settingsRepo *repository.SettingsRepository
}

func NewUserHandler(db *sql.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		profileRepo:  repository.NewProfileRepository(db),
		settingsRepo: repository.NewSettingsRepository(db),
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	profile, err := h.profileRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if profile == nil {
		// Create default profile
		profile = &models.Profile{
			UserID:      userID,
			DisplayName: "",
			Bio:         "",
		}
		if err := h.profileRepo.Create(profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	pkgHttp.Success(w, profile)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var profile models.Profile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile.UserID = userID
	if err := h.profileRepo.Upsert(&profile); err != nil {
		if strings.Contains(err.Error(), "unique_display_name") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"success":false,"error":"nickname_taken"}`))
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, profile)
}

func (h *UserHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	settings, err := h.settingsRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if settings == nil {
		settings = &models.Settings{
			UserID:               userID,
			Language:             "en",
			Theme:                "light",
			NotificationsEnabled: true,
			PrivacyLevel:         "public",
		}
		if err := h.settingsRepo.Create(settings); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	pkgHttp.Success(w, settings)
}

func (h *UserHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var settings models.Settings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	settings.UserID = userID
	if err := h.settingsRepo.Upsert(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, settings)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	profile, err := h.profileRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if profile == nil {
		profile = &models.Profile{UserID: userID}
		if err := h.profileRepo.Create(profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	profile.AvatarURL = req.AvatarURL
	if err := h.profileRepo.Upsert(profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pkgHttp.Success(w, map[string]string{"avatar_url": profile.AvatarURL})
}

func (h *UserHandler) CheckNickname(w http.ResponseWriter, r *http.Request) {
	nickname := r.URL.Query().Get("nickname")
	if nickname == "" {
		pkgHttp.Error(w, http.StatusBadRequest, "nickname parameter is required")
		return
	}

	exists, err := h.profileRepo.NicknameExists(nickname)
	if err != nil {
		pkgHttp.Error(w, http.StatusInternalServerError, "database error")
		return
	}

	pkgHttp.Success(w, map[string]bool{"available": !exists})
}

func (h *UserHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.Success(w, map[string]string{
		"status":  "ok",
		"service": "user-service",
	})
}

func (h *UserHandler) GetPublicProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {
		pkgHttp.Error(w, http.StatusBadRequest, "user ID is required")
		return
	}

	profile, err := h.profileRepo.FindByUserID(userID)
	if err != nil {
		pkgHttp.Error(w, http.StatusInternalServerError, "database error")
		return
	}

	if profile == nil {
		pkgHttp.Error(w, http.StatusNotFound, "user not found")
		return
	}

	publicProfile := map[string]string{
		"id":           profile.UserID,
		"display_name": profile.DisplayName,
		"avatar_url":   profile.AvatarURL,
	}

	pkgHttp.Success(w, publicProfile)
}
