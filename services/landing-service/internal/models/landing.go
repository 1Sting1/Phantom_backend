package models

import "time"

// NavigationItem represents a navigation menu item
type NavigationItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Href  string `json:"href"`
	Type  string `json:"type"` // "link", "button", "dropdown"
	Order int    `json:"order"`
}

// CarouselSlide represents a carousel slide
type CarouselSlide struct {
	ID          string  `json:"id"`
	ImageURL    string  `json:"imageUrl"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Order       int     `json:"order"`
}

// Feature represents a feature card
type Feature struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Icon        *string `json:"icon,omitempty"`
	Order       int     `json:"order"`
	Language    string  `json:"language"` // "ru", "en", etc.
}

// WaitlistEntry represents a waitlist entry
type WaitlistEntry struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Telegram  *string   `json:"telegram,omitempty"`
	Discord   *string   `json:"discord,omitempty"`
	Status    string    `json:"status"` // "pending", "approved", "rejected"
	CreatedAt time.Time `json:"created_at"`
}

// Release represents a software release
type Release struct {
	ID          string    `json:"id"`
	Version     string    `json:"version"`
	OS          string    `json:"os"` // "windows", "linux", "mac"
	DownloadURL string    `json:"downloadUrl"`
	Size        *int64    `json:"size,omitempty"` // in bytes
	Changelog   *string   `json:"changelog,omitempty"`
	IsLatest    bool      `json:"isLatest"`
	CreatedAt   time.Time `json:"created_at"`
}

// Page represents a static page (privacy, terms, etc.)
type Page struct {
	ID        string    `json:"id"`
	Slug      string    `json:"slug"` // "privacy", "terms", etc.
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FooterLink represents a footer link
type FooterLink struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Href     string `json:"href"`
	Category string `json:"category"` // "legal", "social", "support", etc.
	Order    int    `json:"order"`
}

// CommunityLink represents a community link (Discord, Telegram, etc.)
type CommunityLink struct {
	ID        string     `json:"id"`
	URL       string     `json:"url"`
	Type      string     `json:"type"` // "discord", "telegram", "forum"
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	IsActive  bool       `json:"isActive"`
}

// Request/Response models
type WaitlistRequest struct {
	Email    string  `json:"email"`
	Telegram *string `json:"telegram,omitempty"`
	Discord  *string `json:"discord,omitempty"`
}

type WaitlistResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// DownloadToken represents a token for downloading releases
type DownloadToken struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	Email     *string   `json:"email,omitempty"` // Optional: if download requires waitlist
	ExpiresAt time.Time `json:"expiresAt"`
	IsUsed    bool      `json:"isUsed"`
	CreatedAt time.Time `json:"createdAt"`
}

// WaitlistConfirmationToken represents a token for confirming waitlist entry
type WaitlistConfirmationToken struct {
	ID         string    `json:"id"`
	Token      string    `json:"token"`
	WaitlistID string    `json:"waitlistId"`
	Email      string    `json:"email"`
	ExpiresAt  time.Time `json:"expiresAt"`
	IsUsed     bool      `json:"isUsed"`
	CreatedAt  time.Time `json:"createdAt"`
}

// DownloadRequest represents a request to download a release
type DownloadRequest struct {
	OS    string  `json:"os"`              // "windows", "linux", "mac"
	Token *string `json:"token,omitempty"` // Optional: if download requires token
}

// WaitlistConfirmRequest represents a request to confirm waitlist entry
type WaitlistConfirmRequest struct {
	Token string `json:"token"`
}

// DownloadTokenResponse represents a response with download token
type DownloadTokenResponse struct {
	Token       string    `json:"token"`
	ExpiresAt   time.Time `json:"expiresAt"`
	DownloadURL string    `json:"downloadUrl,omitempty"` // Optional: direct download URL
}
