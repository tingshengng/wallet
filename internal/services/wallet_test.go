package services

import (
	"fmt"
	"testing"
	"time"

	"wallet/internal/cache/mock"
	"wallet/internal/models"
	"wallet/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func TestWalletService_Transfer(t *testing.T) {
	mockWalletRepo := &repositories.MockWalletRepository{
		FindByUserIDFunc: func(userID string) (*models.Wallet, error) {
			if userID == "user123" {
				return &models.Wallet{
					ID:        "1",
					UserID:    userID,
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			} else if userID == "user456" {
				return &models.Wallet{
					ID:        "2",
					UserID:    userID,
					Balance:   0.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			}
			return nil, nil
		},
		UpdateFunc: func(wallet *models.Wallet) error {
			return nil
		},
	}

	mockTransactionRepo := &repositories.MockTransactionRepository{
		CreateFunc: func(transaction *models.Transaction) error {
			return nil
		},
	}

	mockCache := &mock.MockCache{
		DeleteFunc: func(key string) {
		},
	}

	walletService := NewWalletService(mockWalletRepo, mockTransactionRepo, mockCache)

	// Test case: Successful transfer
	{
		fromUserID := "user123"
		toUserID := "user456"
		amount := 50.0
		newBalance, _ := walletService.Transfer(fromUserID, toUserID, amount)
		assert.Equal(t, float64(50.0), newBalance) // fromUserID's new balance
	}

	// Test case: Insufficient balance
	{
		fromUserID := "user123"
		toUserID := "user456"
		amount := 150.0
		newBalance, err := walletService.Transfer(fromUserID, toUserID, amount)
		assert.Error(t, err)
		assert.Equal(t, float64(0.0), newBalance)
	}
}

func TestWalletService_Deposit(t *testing.T) {
	mockWalletRepo := &repositories.MockWalletRepository{
		FindByUserIDFunc: func(userID string) (*models.Wallet, error) {
			if userID == "user123" {
				return &models.Wallet{
					ID:        "1",
					UserID:    userID,
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			} else if userID == "user789" {
				return nil, fmt.Errorf("wallet not found")
			}
			return nil, nil
		},
		UpdateFunc: func(wallet *models.Wallet) error {
			return nil
		},
	}

	mockTransactionRepo := &repositories.MockTransactionRepository{
		CreateFunc: func(transaction *models.Transaction) error {
			return nil
		},
	}

	mockCache := &mock.MockCache{
		DeleteFunc: func(key string) {
		},
	}

	walletService := NewWalletService(mockWalletRepo, mockTransactionRepo, mockCache)

	// Test case: Successful deposit
	{
		userID := "user123"
		amount := 100.0
		newBalance, _ := walletService.Deposit(userID, amount)
		assert.Equal(t, float64(200.0), newBalance) // 100.0 (initial) + 100.0 (deposit)
	}

	// Test case: Wallet not found
	{
		userID := "user789"
		amount := 50.0
		newBalance, err := walletService.Deposit(userID, amount)
		assert.Error(t, err)
		assert.Equal(t, float64(0.0), newBalance)
	}
}

func TestWalletService_Withdraw(t *testing.T) {
	mockWalletRepo := &repositories.MockWalletRepository{
		FindByUserIDFunc: func(userID string) (*models.Wallet, error) {
			if userID == "user123" {
				return &models.Wallet{
					ID:        "1",
					UserID:    userID,
					Balance:   100.0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			} else if userID == "user456" {
				return nil, fmt.Errorf("wallet not found")
			}
			return nil, nil
		},
		UpdateFunc: func(wallet *models.Wallet) error {
			return nil
		},
	}

	mockTransactionRepo := &repositories.MockTransactionRepository{
		CreateFunc: func(transaction *models.Transaction) error {
			return nil
		},
	}

	mockCache := &mock.MockCache{
		DeleteFunc: func(key string) {
		},
	}

	walletService := NewWalletService(mockWalletRepo, mockTransactionRepo, mockCache)

	// Test case: Successful withdraw
	{
		userID := "user123"
		amount := 100.0
		newBalance, _ := walletService.Withdraw(userID, amount)
		assert.Equal(t, float64(0.0), newBalance)
	}

	// Test case: Insufficient balance
	{
		userID := "user123"
		amount := 150.0
		_, err := walletService.Withdraw(userID, amount)
		assert.Error(t, err)
	}

	// Test case: Wallet not found
	{
		userID := "user456"
		amount := 50.0
		_, err := walletService.Withdraw(userID, amount)
		assert.Error(t, err)
	}
}
