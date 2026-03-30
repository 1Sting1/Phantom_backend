package repository

import (
	"database/sql"

	"Phantom_backend/services/landing-service/internal/models"
)

type ReleaseRepository struct {
	db *sql.DB
}

func NewReleaseRepository(db *sql.DB) *ReleaseRepository {
	return &ReleaseRepository{db: db}
}

func (r *ReleaseRepository) FindLatest(os string) (*models.Release, error) {
	release := &models.Release{}
	query := `SELECT id, version, os, download_url, size, changelog, is_latest, created_at 
	          FROM releases WHERE os = $1 AND is_latest = TRUE ORDER BY created_at DESC LIMIT 1`

	var size sql.NullInt64
	var changelog sql.NullString
	err := r.db.QueryRow(query, os).Scan(
		&release.ID, &release.Version, &release.OS, &release.DownloadURL,
		&size, &changelog, &release.IsLatest, &release.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if size.Valid {
		release.Size = &size.Int64
	}
	if changelog.Valid {
		release.Changelog = &changelog.String
	}

	return release, nil
}
