package config

import (
	"database/sql"
	"os"
	"strings"

	"Phantom_backend/pkg/database"
)

type Config struct {
	Port      string
	JWTSecret string
	DB        *sql.DB
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8002"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
	}
}

func (c *Config) InitDB() (*sql.DB, error) {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "user_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		return nil, err
	}

	c.DB = db
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// RunMigrations applies user/profile schema (idempotent).
func RunMigrations(db *sql.DB) error {
	for _, block := range migrationBlocks {
		for _, s := range splitSQL(block) {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			if _, err := db.Exec(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func splitSQL(block string) []string {
	var out []string
	for _, s := range strings.Split(block, ";") {
		if t := strings.TrimSpace(s); t != "" {
			out = append(out, t+";")
		}
	}
	return out
}

var migrationBlocks = []string{
	`CREATE TABLE IF NOT EXISTS profiles (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    bio TEXT,
    avatar_url VARCHAR(500),
    country VARCHAR(100),
    timezone VARCHAR(100),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_profiles_user_id ON profiles(user_id);
CREATE TABLE IF NOT EXISTS user_settings (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    theme VARCHAR(20) DEFAULT 'light',
    notifications_enabled BOOLEAN DEFAULT TRUE,
    privacy_level VARCHAR(20) DEFAULT 'public'
);
CREATE INDEX IF NOT EXISTS idx_user_settings_user_id ON user_settings(user_id);`,
	`CREATE UNIQUE INDEX IF NOT EXISTS unique_display_name ON profiles (display_name) WHERE display_name != '';`,
}
