package repository

import (
	"database/sql"

	"Phantom_backend/services/forum-service/internal/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description, slug, parent_id, \"order\", icon, created_at FROM categories ORDER BY \"order\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Slug, &c.ParentID, &c.Order, &c.Icon, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
