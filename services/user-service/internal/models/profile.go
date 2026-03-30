package models

import "time"

type Profile struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	Country     string    `json:"country"`
	Timezone    string    `json:"timezone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Settings struct {
	ID                   string `json:"id"`
	UserID               string `json:"user_id"`
	Language             string `json:"language"`
	Theme                string `json:"theme"`
	NotificationsEnabled bool   `json:"notifications_enabled"`
	PrivacyLevel         string `json:"privacy_level"`
}
