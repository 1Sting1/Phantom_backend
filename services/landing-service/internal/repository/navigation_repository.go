package repository

import (
	"database/sql"

	"Phantom_backend/services/landing-service/internal/models"
)

type NavigationRepository struct {
	db *sql.DB
}

func NewNavigationRepository(db *sql.DB) *NavigationRepository {
	return &NavigationRepository{db: db}
}

func (r *NavigationRepository) FindAll() ([]models.NavigationItem, error) {
	rows, err := r.db.Query("SELECT id, title, href, type, \"order\" FROM navigation_items ORDER BY \"order\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.NavigationItem
	for rows.Next() {
		var item models.NavigationItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Href, &item.Type, &item.Order); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
