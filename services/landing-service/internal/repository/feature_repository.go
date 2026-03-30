package repository

import (
	"database/sql"

	"Phantom_backend/services/landing-service/internal/models"
)

type FeatureRepository struct {
	db *sql.DB
}

func NewFeatureRepository(db *sql.DB) *FeatureRepository {
	return &FeatureRepository{db: db}
}

func (r *FeatureRepository) FindAll(language string) ([]models.Feature, error) {
	query := "SELECT id, title, description, icon, \"order\", language FROM features WHERE language = $1 ORDER BY \"order\""
	rows, err := r.db.Query(query, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var features []models.Feature
	for rows.Next() {
		var feature models.Feature
		var icon sql.NullString
		if err := rows.Scan(&feature.ID, &feature.Title, &feature.Description, &icon, &feature.Order, &feature.Language); err != nil {
			return nil, err
		}
		if icon.Valid {
			feature.Icon = &icon.String
		}
		features = append(features, feature)
	}

	return features, nil
}
