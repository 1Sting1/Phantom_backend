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
	"Phantom_backend/services/landing-service/internal/config"
	"Phantom_backend/services/landing-service/internal/handlers"
	"Phantom_backend/services/landing-service/internal/middleware"

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

	db, err := cfg.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	router := mux.NewRouter()

	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.Recovery(log))

	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	landingHandler := handlers.NewLandingHandler(db)

	// Public API routes
	publicAPI := router.PathPrefix("/api/v1/public").Subrouter()
	publicAPI.HandleFunc("/navigation", landingHandler.GetNavigation).Methods("GET")
	publicAPI.HandleFunc("/carousel", landingHandler.GetCarousel).Methods("GET")
	publicAPI.HandleFunc("/features", landingHandler.GetFeatures).Methods("GET")
	publicAPI.HandleFunc("/waitlist", landingHandler.AddToWaitlist).Methods("POST")
	publicAPI.HandleFunc("/waitlist/confirm", landingHandler.ConfirmWaitlist).Methods("POST")
	publicAPI.HandleFunc("/releases/latest", landingHandler.GetLatestRelease).Methods("GET")
	publicAPI.HandleFunc("/download", landingHandler.Download).Methods("GET")
	publicAPI.HandleFunc("/download-token", landingHandler.GetDownloadToken).Methods("GET")
	publicAPI.HandleFunc("/pages", landingHandler.GetPage).Methods("GET")
	publicAPI.HandleFunc("/footer-links", landingHandler.GetFooterLinks).Methods("GET")
	publicAPI.HandleFunc("/community-link", landingHandler.GetCommunityLink).Methods("GET")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Info("Starting Landing Service", zap.String("port", cfg.Port))
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
