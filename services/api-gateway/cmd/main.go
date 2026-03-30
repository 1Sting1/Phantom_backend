package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Phantom_backend/pkg/logger"
	"Phantom_backend/services/api-gateway/internal/config"
	"Phantom_backend/services/api-gateway/internal/handlers"
	"Phantom_backend/services/api-gateway/internal/middleware"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log, err := logger.NewLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.Recovery(log))

	// Health check
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.RateLimiter())

	// Auth routes
	authRouter := api.PathPrefix("/auth").Subrouter()
	handlers.RegisterAuthRoutes(authRouter, cfg)

	// Me route (requires auth) - proxy to auth service /me endpoint
	meRouter := api.PathPrefix("/me").Subrouter()
	meRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	meRouter.HandleFunc("", handlers.ProxyRequest(cfg, cfg.AuthServiceURL)).Methods("GET")

	// User routes
	userRouter := api.PathPrefix("/user").Subrouter()
	userRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	handlers.RegisterUserRoutes(userRouter, cfg)

	// Public User routes
	publicUserRouter := api.PathPrefix("/public/user").Subrouter()
	publicUserRouter.HandleFunc("/check-nickname", handlers.ProxyRequest(cfg, cfg.UserServiceURL)).Methods("GET")
	publicUserRouter.HandleFunc("/profile/{id}", handlers.ProxyRequest(cfg, cfg.UserServiceURL)).Methods("GET")

	// Forum routes (optional auth: X-User-ID set when Bearer token present for create/edit/delete)
	forumRouter := api.PathPrefix("/forum").Subrouter()
	forumRouter.Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))
	handlers.RegisterForumRoutes(forumRouter, cfg)

	// Shop routes (optional auth: order attribution when Bearer present)
	shopRouter := api.PathPrefix("/shop").Subrouter()
	shopRouter.Use(middleware.OptionalAuthMiddleware(cfg.JWTSecret))
	handlers.RegisterShopRoutes(shopRouter, cfg)

	// Localization routes
	locRouter := api.PathPrefix("/localization").Subrouter()
	handlers.RegisterLocalizationRoutes(locRouter, cfg)

	// Notifications routes
	notifRouter := api.PathPrefix("/notifications").Subrouter()
	notifRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	handlers.RegisterNotificationRoutes(notifRouter, cfg)

	// Public landing routes (no auth required)
	publicRouter := api.PathPrefix("/public").Subrouter()
	handlers.RegisterLandingRoutes(publicRouter, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Info("Starting API Gateway", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited")
}
