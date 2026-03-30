package models

import "time"

type Sticker struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PreviewURL  string    `json:"preview_url"`
	FileURL     string    `json:"file_url"`
	Price       float64   `json:"price"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
}

type Pack struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PreviewURL  string    `json:"preview_url"`
	Price       float64   `json:"price"`
	Discount    float64   `json:"discount"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
