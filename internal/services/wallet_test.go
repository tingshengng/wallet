package services

import (
	"database/sql"
	"testing"

	cachemock "wallet/internal/cache/mock"
	"wallet/internal/models"
	"wallet/internal/repositories"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTests initializes a mock DB and repositories for testing
func setupTests(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *repositories.MockWalletRepository, *repositories.MockTransactionRepository, *cachemock.MockCache, WalletService) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
	}

	mockWalletRepo := &repositories.MockWalletRepository{}
	mockTransactionRepo := &repositories.MockTransactionRepository{}
	mockCache := &cachemock.MockCache{}

	// Mock the DB transaction methods
	mockWalletRepo.DBFunc = func() *gorm.DB {
		return gormDB
	}
	mockWalletRepo.WithTxFunc = func(tx interface{}) repositories.WalletRepository {
		return mockWalletRepo // Return the same mock, as we control its behavior
	}
	mockTransactionRepo.WithTxFunc = func(tx interface{}) repositories.TransactionRepository {
		return mockTransactionRepo // Return the same mock
	}

	walletService := NewWalletService(mockWalletRepo, mockTransactionRepo, mockCache)

	return db, mock, mockWalletRepo, mockTransactionRepo, mockCache, walletService
}

func TestWalletService_Deposit(t *testing.T) {
	t.Run("successful deposit", func(t *testing.T) {
		db, mock, mockWalletRepo, mockTransactionRepo, mockCache, walletService := setupTests(t)
		defer db.Close()

		userID := "user123"
		amount := 100.0
		initialBalance := 50.0

		mockWalletRepo.FindByUserIDFunc = func(uid string) (*models.Wallet, error) {
			assert.Equal(t, userID, uid)
			return &models.Wallet{ID: "wallet1", UserID: userID, Balance: initialBalance}, nil
		}

		mock.ExpectBegin()

		mockTransactionRepo.CreateFunc = func(tx *models.Transaction) error {
			assert.Equal(t, userID, tx.FromUserID)
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, models.TransactionTypeDeposit, tx.Type)
			return nil
		}

		mockWalletRepo.UpdateFunc = func(w *models.Wallet) error {
			assert.Equal(t, userID, w.UserID)
			assert.Equal(t, initialBalance+amount, w.Balance)
			return nil
		}

		mock.ExpectCommit()

		mockCache.DeleteFunc = func(key string) {
			assert.Equal(t, userID, key)
		}

		newBalance, err := walletService.Deposit(userID, amount)

		assert.Nil(t, err)
		assert.Equal(t, initialBalance+amount, newBalance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWalletService_Withdraw(t *testing.T) {
	t.Run("successful withdraw", func(t *testing.T) {
		db, mock, mockWalletRepo, mockTransactionRepo, mockCache, walletService := setupTests(t)
		defer db.Close()

		userID := "user123"
		amount := 50.0
		initialBalance := 100.0

		mockWalletRepo.FindByUserIDFunc = func(uid string) (*models.Wallet, error) {
			assert.Equal(t, userID, uid)
			return &models.Wallet{ID: "wallet1", UserID: userID, Balance: initialBalance}, nil
		}

		mock.ExpectBegin()

		mockTransactionRepo.CreateFunc = func(tx *models.Transaction) error {
			assert.Equal(t, userID, tx.FromUserID)
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, models.TransactionTypeWithdraw, tx.Type)
			return nil
		}

		mockWalletRepo.UpdateFunc = func(w *models.Wallet) error {
			assert.Equal(t, userID, w.UserID)
			assert.Equal(t, initialBalance-amount, w.Balance)
			return nil
		}

		mock.ExpectCommit()

		mockCache.DeleteFunc = func(key string) {
			assert.Equal(t, userID, key)
		}

		newBalance, err := walletService.Withdraw(userID, amount)

		assert.Nil(t, err)
		assert.Equal(t, initialBalance-amount, newBalance)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insufficient balance", func(t *testing.T) {
		db, _, mockWalletRepo, _, _, walletService := setupTests(t)
		defer db.Close()

		userID := "user123"
		amount := 150.0
		initialBalance := 100.0

		mockWalletRepo.FindByUserIDFunc = func(uid string) (*models.Wallet, error) {
			return &models.Wallet{ID: "wallet1", UserID: userID, Balance: initialBalance}, nil
		}

		_, apiErr := walletService.Withdraw(userID, amount)

		assert.Error(t, apiErr)
		assert.Equal(t, "Insufficient balance", apiErr.Message)
	})
}

func TestWalletService_Transfer(t *testing.T) {
	t.Run("successful transfer", func(t *testing.T) {
		db, mock, mockWalletRepo, mockTransactionRepo, mockCache, walletService := setupTests(t)
		defer db.Close()

		fromUserID := "user123"
		toUserID := "user456"
		amount := 50.0
		fromInitialBalance := 100.0
		toInitialBalance := 20.0

		mockWalletRepo.FindByUserIDFunc = func(userID string) (*models.Wallet, error) {
			if userID == fromUserID {
				return &models.Wallet{ID: "wallet1", UserID: fromUserID, Balance: fromInitialBalance}, nil
			}
			if userID == toUserID {
				return &models.Wallet{ID: "wallet2", UserID: toUserID, Balance: toInitialBalance}, nil
			}
			return nil, nil
		}

		mock.ExpectBegin()

		mockTransactionRepo.CreateFunc = func(tx *models.Transaction) error {
			assert.Equal(t, fromUserID, tx.FromUserID)
			assert.Equal(t, toUserID, tx.ToUserID)
			assert.Equal(t, amount, tx.Amount)
			assert.Equal(t, models.TransactionTypeTransfer, tx.Type)
			return nil
		}

		updateCount := 0
		mockWalletRepo.UpdateFunc = func(w *models.Wallet) error {
			updateCount++
			if w.UserID == fromUserID {
				assert.Equal(t, fromInitialBalance-amount, w.Balance)
			} else {
				assert.Equal(t, toInitialBalance+amount, w.Balance)
			}
			return nil
		}

		mock.ExpectCommit()

		deleteCount := 0
		mockCache.DeleteFunc = func(key string) {
			deleteCount++
			assert.Contains(t, []string{fromUserID, toUserID}, key)
		}

		newBalance, err := walletService.Transfer(fromUserID, toUserID, amount)

		assert.Nil(t, err)
		assert.Equal(t, fromInitialBalance-amount, newBalance)
		assert.Equal(t, 2, updateCount)
		assert.Equal(t, 2, deleteCount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("insufficient balance", func(t *testing.T) {
		db, _, mockWalletRepo, _, _, walletService := setupTests(t)
		defer db.Close()

		fromUserID := "user123"
		toUserID := "user456"
		amount := 150.0
		fromInitialBalance := 100.0

		mockWalletRepo.FindByUserIDFunc = func(userID string) (*models.Wallet, error) {
			if userID == fromUserID {
				return &models.Wallet{ID: "wallet1", UserID: fromUserID, Balance: fromInitialBalance}, nil
			}
			if userID == toUserID {
				return &models.Wallet{ID: "wallet2", UserID: toUserID, Balance: 0}, nil
			}
			return nil, nil
		}

		_, apiErr := walletService.Transfer(fromUserID, toUserID, amount)

		assert.Error(t, apiErr)
		assert.Equal(t, "Insufficient balance", apiErr.Message)
	})
}
