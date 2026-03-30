package handlers

import (
	"net/http"

	pkgHttp "Phantom_backend/pkg/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	pkgHttp.Success(w, map[string]string{
		"status":  "ok",
		"service": "api-gateway",
	})
}
