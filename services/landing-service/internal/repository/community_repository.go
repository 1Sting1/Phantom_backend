package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/landing-service/internal/models"
)

type CommunityRepository struct {
	db *sql.DB
}

func NewCommunityRepository(db *sql.DB) *CommunityRepository {
	return &CommunityRepository{db: db}
}

func (r *CommunityRepository) FindActive() (*models.CommunityLink, error) {
	link := &models.CommunityLink{}
	query := `SELECT id, url, type, expires_at, is_active 
	          FROM community_links WHERE is_active = TRUE 
	          AND (expires_at IS NULL OR expires_at > $1) 
	          ORDER BY created_at DESC LIMIT 1`

	var expiresAt sql.NullTime
	err := r.db.QueryRow(query, time.Now()).Scan(
		&link.ID, &link.URL, &link.Type, &expiresAt, &link.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if expiresAt.Valid {
		link.ExpiresAt = &expiresAt.Time
	}

	return link, nil
}
