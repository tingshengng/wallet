package models

import (
	"time"
)

type UserToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id" gorm:"index:idx_user_token_user_id"`
	Token     string    `json:"token" gorm:"index:idx_user_token_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
