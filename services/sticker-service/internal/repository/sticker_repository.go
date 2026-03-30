package repository

import (
	"database/sql"
	"time"

	"Phantom_backend/services/sticker-service/internal/models"
)

type StickerRepository struct {
	db *sql.DB
}

func NewStickerRepository(db *sql.DB) *StickerRepository {
	return &StickerRepository{db: db}
}

func (r *StickerRepository) FindAll() ([]models.Sticker, error) {
	rows, err := r.db.Query("SELECT id, name, description, preview_url, file_url, price, category_id, created_at, is_active FROM stickers WHERE is_active = true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stickers []models.Sticker
	for rows.Next() {
		var s models.Sticker
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.PreviewURL, &s.FileURL,
			&s.Price, &s.CategoryID, &s.CreatedAt, &s.IsActive); err != nil {
			return nil, err
		}
		stickers = append(stickers, s)
	}

	return stickers, nil
}

// OrderRepository persists shop orders (table shop_orders, created by config.RunMigrations).
type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *models.Order) error {
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	_, err := r.db.Exec(
		`INSERT INTO shop_orders (id, user_id, total_amount, status, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		order.ID, order.UserID, order.TotalAmount, order.Status, order.CreatedAt,
	)
	return err
}
