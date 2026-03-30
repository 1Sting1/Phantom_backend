package repository

import (
	"database/sql"

	"Phantom_backend/services/user-service/internal/models"
	"github.com/google/uuid"
)

type SettingsRepository struct {
	db *sql.DB
}

func NewSettingsRepository(db *sql.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r *SettingsRepository) Create(settings *models.Settings) error {
	settings.ID = uuid.New().String()

	query := `INSERT INTO user_settings (id, user_id, language, theme, notifications_enabled, privacy_level) 
	          VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, settings.ID, settings.UserID, settings.Language,
		settings.Theme, settings.NotificationsEnabled, settings.PrivacyLevel)
	return err
}

func (r *SettingsRepository) FindByUserID(userID string) (*models.Settings, error) {
	settings := &models.Settings{}
	query := `SELECT id, user_id, language, theme, notifications_enabled, privacy_level 
	          FROM user_settings WHERE user_id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&settings.ID, &settings.UserID, &settings.Language, &settings.Theme,
		&settings.NotificationsEnabled, &settings.PrivacyLevel,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return settings, err
}

func (r *SettingsRepository) Upsert(settings *models.Settings) error {
	existing, _ := r.FindByUserID(settings.UserID)

	if existing == nil {
		return r.Create(settings)
	}

	settings.ID = existing.ID

	query := `UPDATE user_settings SET language = $1, theme = $2, notifications_enabled = $3, privacy_level = $4 
	          WHERE user_id = $5`

	_, err := r.db.Exec(query, settings.Language, settings.Theme,
		settings.NotificationsEnabled, settings.PrivacyLevel, settings.UserID)
	return err
}
