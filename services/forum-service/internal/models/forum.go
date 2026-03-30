package models

import "time"

// GuestUserID is used for forum posts when no JWT is sent (anonymous replies).
const GuestUserID = "00000000-0000-0000-0000-000000000000"

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	Slug        string    `json:"slug"`
	ParentID    *string   `json:"parent_id"`
	Order       int       `json:"order"`
	Icon        *string   `json:"icon"`
	CreatedAt   time.Time `json:"created_at"`
}

type Thread struct {
	ID         string    `json:"id"`
	CategoryID string    `json:"category_id"`
	UserID     string    `json:"user_id"`
	Title      string    `json:"title"`
	Slug       string    `json:"slug"`
	Content    string    `json:"content"`
	IsPinned   bool      `json:"is_pinned"`
	IsLocked   bool      `json:"is_locked"`
	ViewsCount int       `json:"views_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Post struct {
	ID           string    `json:"id"`
	ThreadID     string    `json:"thread_id"`
	UserID       string    `json:"user_id"`
	ParentPostID string    `json:"parent_post_id"`
	Content      string    `json:"content"`
	IsEdited     bool      `json:"is_edited"`
	EditedAt     time.Time `json:"edited_at"`
	CreatedAt    time.Time `json:"created_at"`
}
