package models

import (
	"time"
)

type Wallet struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id" gorm:"index:idx_wallet_user_id"`
	Balance   float64    `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
