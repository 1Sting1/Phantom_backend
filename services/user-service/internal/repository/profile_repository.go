package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/user-service/internal/models"
	"github.com/google/uuid"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Create(profile *models.Profile) error {
	profile.ID = uuid.New().String()
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	query := `INSERT INTO profiles (id, user_id, display_name, bio, avatar_url, country, timezone, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, profile.ID, profile.UserID, profile.DisplayName, profile.Bio,
		profile.AvatarURL, profile.Country, profile.Timezone, profile.CreatedAt, profile.UpdatedAt)
	return err
}

func (r *ProfileRepository) FindByUserID(userID string) (*models.Profile, error) {
	profile := &models.Profile{}
	query := `SELECT id, user_id, display_name, bio, avatar_url, country, timezone, created_at, updated_at 
	          FROM profiles WHERE user_id = $1`

	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.DisplayName, &profile.Bio,
		&profile.AvatarURL, &profile.Country, &profile.Timezone, &profile.CreatedAt, &profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return profile, err
}

func (r *ProfileRepository) Upsert(profile *models.Profile) error {
	existing, _ := r.FindByUserID(profile.UserID)

	if existing == nil {
		return r.Create(profile)
	}

	profile.ID = existing.ID
	profile.UpdatedAt = time.Now()

	query := `UPDATE profiles SET display_name = $1, bio = $2, avatar_url = $3, country = $4, timezone = $5, updated_at = $6 
	          WHERE user_id = $7`

	_, err := r.db.Exec(query, profile.DisplayName, profile.Bio, profile.AvatarURL,
		profile.Country, profile.Timezone, profile.UpdatedAt, profile.UserID)
	return err
}

func (r *ProfileRepository) NicknameExists(nickname string) (bool, error) {
	var id string
	query := `SELECT id FROM profiles WHERE display_name = $1 LIMIT 1`
	err := r.db.QueryRow(query, nickname).Scan(&id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
