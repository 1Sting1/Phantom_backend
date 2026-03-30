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
	"Phantom_backend/services/notific
	"Phantom_backend/services/notification-service/internal/config"
	"Phantom_backend/services/notification-service/internal/handlers"
	"Phantom_backend/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log, err := logger.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger", zap.Error(err))
	}
	defer log.Sync()

	router := mux.NewRouter()

	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.Recovery(log))

	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	notifHandler := handlers.NewNotificationHandler()
	api := router.PathPrefix("/api/v1/notifications").Subrouter()
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	api.HandleFunc("", notifHandler.GetNotifications).Methods("GET")
	api.HandleFunc("/{id}/read", notifHandler.MarkAsRead).Methods("PUT")
	api.HandleFunc("/preferences", notifHandler.GetPreferences).Methods("GET")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Info("Starting Notification Service", zap.String("port", cfg.Port))
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
d
}