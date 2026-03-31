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
		Port: getEnv("PORT", "8007"),
	}
}

func (c *Config) InitDB() (*sql.DB, error) {
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "landing_db"),
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

// RunMigrations applies landing schema on startup so services can boot without manual migration step.
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
	`CREATE TABLE IF NOT EXISTS navigation_items (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    href VARCHAR(500) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'link',
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS carousel_slides (
    id VARCHAR(255) PRIMARY KEY,
    image_url VARCHAR(1000) NOT NULL,
    title VARCHAR(500),
    description TEXT,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS features (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    icon VARCHAR(255),
    "order" INTEGER DEFAULT 0,
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS waitlist_entries (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    telegram VARCHAR(255),
    discord VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS releases (
    id VARCHAR(255) PRIMARY KEY,
    version VARCHAR(50) NOT NULL,
    os VARCHAR(50) NOT NULL,
    download_url VARCHAR(1000) NOT NULL,
    size BIGINT,
    changelog TEXT,
    is_latest BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS pages (
    id VARCHAR(255) PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS footer_links (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    href VARCHAR(500) NOT NULL,
    category VARCHAR(50) NOT NULL,
    "order" INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS community_links (
    id VARCHAR(255) PRIMARY KEY,
    url VARCHAR(1000) NOT NULL,
    type VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS download_tokens (
    id VARCHAR(255) PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS waitlist_confirmation_tokens (
    id VARCHAR(255) PRIMARY KEY,
    token VARCHAR(255) UNIQUE NOT NULL,
    waitlist_id VARCHAR(255) NOT NULL REFERENCES waitlist_entries(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE waitlist_entries ADD COLUMN IF NOT EXISTS confirmation_token_id VARCHAR(255);
ALTER TABLE waitlist_entries ADD COLUMN IF NOT EXISTS confirmed_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_navigation_order ON navigation_items("order");
CREATE INDEX IF NOT EXISTS idx_carousel_order ON carousel_slides("order");
CREATE INDEX IF NOT EXISTS idx_features_order ON features("order");
CREATE INDEX IF NOT EXISTS idx_features_language ON features(language);
CREATE INDEX IF NOT EXISTS idx_waitlist_email ON waitlist_entries(email);
CREATE INDEX IF NOT EXISTS idx_waitlist_status ON waitlist_entries(status);
CREATE INDEX IF NOT EXISTS idx_releases_os ON releases(os);
CREATE INDEX IF NOT EXISTS idx_releases_latest ON releases(is_latest) WHERE is_latest = TRUE;
CREATE INDEX IF NOT EXISTS idx_pages_slug ON pages(slug);
CREATE INDEX IF NOT EXISTS idx_pages_language ON pages(language);
CREATE INDEX IF NOT EXISTS idx_footer_category ON footer_links(category);
CREATE INDEX IF NOT EXISTS idx_community_active ON community_links(is_active) WHERE is_active = TRUE;
CREATE INDEX IF NOT EXISTS idx_download_tokens_token ON download_tokens(token);
CREATE INDEX IF NOT EXISTS idx_download_tokens_email ON download_tokens(email);
CREATE INDEX IF NOT EXISTS idx_download_tokens_expires ON download_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_waitlist_confirmation_tokens_token ON waitlist_confirmation_tokens(token);
CREATE INDEX IF NOT EXISTS idx_waitlist_confirmation_tokens_waitlist_id ON waitlist_confirmation_tokens(waitlist_id);`,
}
