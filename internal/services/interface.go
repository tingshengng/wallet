package services

import (
	"wallet/internal/models"
)

type WalletService interface {
	Deposit(userID string, amount float64) (float64, *APIError)
	Withdraw(userID string, amount float64) (float64, *APIError)
	Transfer(fromUserID, toUserID string, amount float64) (float64, *APIError)
	GetBalance(userID string) (float64, *APIError)
	GetTransactionHistory(userID string, page, pageSize int, transactionType, status string) ([]models.Transaction, *APIError)
}
