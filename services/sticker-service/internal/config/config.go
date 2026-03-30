package config

import (
	"database/sql"
	"os"
	"strings"

	"Phantom_backend/pkg/database"
)

type Config struct {
	Port string
	DB   *sql.DB
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8004"),
	}
}

func (c *Config) InitDB() (*sql.DB, error) {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "shop_db"),
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

// RunMigrations applies shop schema (idempotent).
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
	`CREATE TABLE IF NOT EXISTS stickers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    preview_url VARCHAR(500),
    file_url VARCHAR(500) NOT NULL,
    price DECIMAL(10,2) DEFAULT 0,
    category_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS sticker_packs (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    preview_url VARCHAR(500),
    price DECIMAL(10,2) DEFAULT 0,
    discount DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL
);`,
	`CREATE TABLE IF NOT EXISTS shop_orders (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_shop_orders_user_id ON shop_orders(user_id);`,
}
