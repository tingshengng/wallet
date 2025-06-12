package models

import (
	"time"
)

type Transaction struct {
	ID         string     `json:"id"`
	FromUserID string     `json:"from_user_id" gorm:"index:idx_transaction_from_user_id"`
	ToUserID   string     `json:"to_user_id" gorm:"index:idx_transaction_to_user_id"`
	Amount     float64    `json:"amount"`
	Type       string     `json:"type"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

const (
	TransactionTypeDeposit   = "deposit"
	TransactionTypeWithdraw  = "withdraw"
	TransactionTypeTransfer  = "transfer"
	TransactionStatusPending = "pending"
	TransactionStatusSuccess = "success"
	TransactionStatusFailed  = "failed"
)
