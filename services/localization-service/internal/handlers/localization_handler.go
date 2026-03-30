package handlers

import (
	"net/http"

	httputil "Phantom_backend/pkg/http"
	"github.com/gorilla/mux"
)

type LocalizationHandler struct {
	translations map[string]map[string]string
}

func NewLocalizationHandler() *LocalizationHandler {
	translations := map[string]map[string]string{
		"en": {
			"welcome":  "Welcome",
			"login":    "Login",
			"register": "Register",
		},
		"ru": {
			"welcome":  "Добро пожаловать",
			"login":    "Войти",
			"register": "Регистрация",
		},
	}

	return &LocalizationHandler{
		translations: translations,
	}
}

func (h *LocalizationHandler) GetTranslations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	language := vars["language"]

	translations, exists := h.translations[language]
	if !exists {
		translations = h.translations["en"] // fallback to English
	}

	httputil.Success(w, translations)
}

func (h *LocalizationHandler) GetLanguages(w http.ResponseWriter, r *http.Request) {
	languages := []map[string]string{
		{"code": "en", "name": "English", "native_name": "English"},
		{"code": "ru", "name": "Russian", "native_name": "Русский"},
	}

	httputil.Success(w, languages)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, map[string]string{
		"status":  "ok",
		"service": "localization-service",
	})
}
