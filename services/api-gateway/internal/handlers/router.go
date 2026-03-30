package handlers

import (
	"Phantom_backend/services/api-gateway/internal/config"

	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/register", ProxyRequest(cfg, cfg.AuthServiceURL)).Methods("POST")
	router.HandleFunc("/login", ProxyRequest(cfg, cfg.AuthServiceURL)).Methods("POST")
	router.HandleFunc("/logout", ProxyRequest(cfg, cfg.AuthServiceURL)).Methods("POST")
	router.HandleFunc("/refresh", ProxyRequest(cfg, cfg.AuthServiceURL)).Methods("POST")
	router.HandleFunc("/me", ProxyRequest(cfg, cfg.AuthServiceURL)).Methods("GET")
}

func RegisterUserRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/profile", ProxyRequest(cfg, cfg.UserServiceURL)).Methods("GET", "PUT")
	router.HandleFunc("/settings", ProxyRequest(cfg, cfg.UserServiceURL)).Methods("GET", "PUT")
	router.HandleFunc("/avatar", ProxyRequest(cfg, cfg.UserServiceURL)).Methods("POST")
}

func RegisterForumRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/categories", ProxyRequest(cfg, cfg.ForumServiceURL)).Methods("GET")
	router.HandleFunc("/threads", ProxyRequest(cfg, cfg.ForumServiceURL)).Methods("GET", "POST")
	router.HandleFunc("/threads/{id}", ProxyRequest(cfg, cfg.ForumServiceURL)).Methods("GET", "PUT", "DELETE")
	router.HandleFunc("/threads/{id}/posts", ProxyRequest(cfg, cfg.ForumServiceURL)).Methods("GET")
	router.HandleFunc("/posts", ProxyRequest(cfg, cfg.ForumServiceURL)).Methods("POST")
	router.HandleFunc("/posts/{id}", ProxyRequest(cfg, cfg.ForumServiceURL)).Methods("PUT", "DELETE")
}

func RegisterShopRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/stickers", ProxyRequest(cfg, cfg.StickerServiceURL)).Methods("GET")
	router.HandleFunc("/packs", ProxyRequest(cfg, cfg.StickerServiceURL)).Methods("GET")
	router.HandleFunc("/order", ProxyRequest(cfg, cfg.StickerServiceURL)).Methods("POST")
}

func RegisterLocalizationRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/{language}", ProxyRequest(cfg, cfg.LocalizationURL)).Methods("GET")
	router.HandleFunc("/languages", ProxyRequest(cfg, cfg.LocalizationURL)).Methods("GET")
}

func RegisterNotificationRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("", ProxyRequest(cfg, cfg.NotificationURL)).Methods("GET")
	router.HandleFunc("/{id}/read", ProxyRequest(cfg, cfg.NotificationURL)).Methods("PUT")
	router.HandleFunc("/preferences", ProxyRequest(cfg, cfg.NotificationURL)).Methods("GET")
}

func RegisterLandingRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/navigation", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/carousel", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/features", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/waitlist", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("POST")
	router.HandleFunc("/waitlist/confirm", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("POST")
	router.HandleFunc("/releases/latest", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/download", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/download-token", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/pages", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/footer-links", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
	router.HandleFunc("/community-link", ProxyRequest(cfg, cfg.LandingServiceURL)).Methods("GET")
}
