package config

import (
	"os"
)

type Config struct {
	Port              string
	JWTSecret         string
	AuthServiceURL    string
	UserServiceURL    string
	ForumServiceURL   string
	StickerServiceURL string
	LocalizationURL   string
	NotificationURL   string
	LandingServiceURL string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8080"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://auth-service:8001"),
		UserServiceURL:    getEnv("USER_SERVICE_URL", "http://user-service:8002"),
		ForumServiceURL:   getEnv("FORUM_SERVICE_URL", "http://forum-service:8003"),
		StickerServiceURL: getEnv("STICKER_SERVICE_URL", "http://sticker-service:8004"),
		LocalizationURL:   getEnv("LOCALIZATION_SERVICE_URL", "http://localization-service:8005"),
		NotificationURL:   getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8006"),
		LandingServiceURL: getEnv("LANDING_SERVICE_URL", "http://landing-service:8007"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
