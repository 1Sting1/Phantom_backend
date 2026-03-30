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
	"Phantom_backend/services/forum-service/internal/config"
	"Phantom_backend/services/forum-service/internal/handlers"
	"Phantom_backend/services/forum-service/internal/middleware"

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

	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	router := mux.NewRouter()

	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())
	router.Use(middleware.Recovery(log))

	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	forumHandler := handlers.NewForumHandler(db)
	api := router.PathPrefix("/api/v1/forum").Subrouter()
	api.HandleFunc("/categories", forumHandler.GetCategories).Methods("GET")
	api.HandleFunc("/threads", forumHandler.GetThreads).Methods("GET")
	api.HandleFunc("/threads", forumHandler.CreateThread).Methods("POST")
	api.HandleFunc("/threads/{id}", forumHandler.GetThread).Methods("GET")
	api.HandleFunc("/threads/{id}", forumHandler.UpdateThread).Methods("PUT")
	api.HandleFunc("/threads/{id}", forumHandler.DeleteThread).Methods("DELETE")
	api.HandleFunc("/threads/{id}/posts", forumHandler.GetThreadPosts).Methods("GET")
	api.HandleFunc("/posts", forumHandler.CreatePost).Methods("POST")
	api.HandleFunc("/posts/{id}", forumHandler.UpdatePost).Methods("PUT")
	api.HandleFunc("/posts/{id}", forumHandler.DeletePost).Methods("DELETE")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	go func() {
		log.Info("Starting Forum Service", zap.String("port", cfg.Port))
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
