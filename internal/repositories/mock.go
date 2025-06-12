package repositories

import (
	"wallet/internal/models"
)

// MockWalletRepository is a mock implementation of WalletRepository
// Manual mocks for repositories

type MockWalletRepository struct {
	WalletRepository
	FindByUserIDFunc func(userID string) (*models.Wallet, error)
	UpdateFunc       func(wallet *models.Wallet) error
	CreateFunc       func(wallet *models.Wallet) error
	DeleteFunc       func(id string) error
}

func (m *MockWalletRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockWalletRepository) FindByUserID(userID string) (*models.Wallet, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(userID)
	}
	return nil, nil
}

func (m *MockWalletRepository) Update(wallet *models.Wallet) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(wallet)
	}
	return nil
}

func (m *MockWalletRepository) Create(wallet *models.Wallet) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(wallet)
	}
	return nil
}

// MockTransactionRepository is a mock implementation of TransactionRepository
type MockTransactionRepository struct {
	TransactionRepository
	CreateFunc func(transaction *models.Transaction) error
	FindFunc   func(userID string, page, pageSize int, transactionType, status string) ([]models.Transaction, error)
	DeleteFunc func(id string) error
}

func (m *MockTransactionRepository) Create(transaction *models.Transaction) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(transaction)
	}
	return nil
}

func (m *MockTransactionRepository) Find(userID string, page, pageSize int, transactionType, status string) ([]models.Transaction, error) {
	if m.FindFunc != nil {
		return m.FindFunc(userID, page, pageSize, transactionType, status)
	}
	return nil, nil
}

func (m *MockTransactionRepository) Delete(id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
