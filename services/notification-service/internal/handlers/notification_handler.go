package handlers

import (
	"net/http"

	httputil "Phantom_backend/pkg/http"
	"github.com/gorilla/mux"
)

type NotificationHandler struct {
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	notifications := []map[string]interface{}{
		{
			"id":      "1",
			"user_id": userID,
			"type":    "info",
			"title":   "Welcome!",
			"message": "Welcome to our platform",
			"is_read": false,
		},
	}

	httputil.Success(w, notifications)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["id"]

	httputil.Success(w, map[string]string{"message": "Notification marked as read"})
}

func (h *NotificationHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	preferences := map[string]interface{}{
		"user_id":        userID,
		"email_enabled":  true,
		"push_enabled":   true,
		"forum_replies":  true,
		"new_stickers":   true,
		"system_updates": true,
	}

	httputil.Success(w, preferences)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httputil.Success(w, map[string]string{
		"status":  "ok",
		"service": "notification-service",
	})
}
