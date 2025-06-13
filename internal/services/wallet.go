package services

import (
	"time"

	"wallet/internal/cache"
	"wallet/internal/models"
	"wallet/internal/repositories"

	"github.com/google/uuid"
)

type walletService struct {
	WalletRepo      repositories.WalletRepository
	TransactionRepo repositories.TransactionRepository
	Cache           cache.Cache
}

func NewWalletService(
	walletRepo repositories.WalletRepository,
	transactionRepo repositories.TransactionRepository,
	cache cache.Cache,
) WalletService {
	return &walletService{
		WalletRepo:      walletRepo,
		TransactionRepo: transactionRepo,
		Cache:           cache,
	}
}

func (s *walletService) Deposit(userID string, amount float64) (float64, *APIError) {
	if amount <= 0 {
		return 0, NewBadRequestError("Invalid amount")
	}

	wallet, err := s.WalletRepo.FindByUserID(userID)
	if err != nil {
		return 0, NewInternalServerError("Failed to get wallet")
	}

	// Start a database transaction for the create and update operations
	tx := s.WalletRepo.DB().Begin()
	if tx.Error != nil {
		return 0, NewInternalServerError("Failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create repository instances with transaction
	walletRepo := s.WalletRepo.WithTx(tx)
	transactionRepo := s.TransactionRepo.WithTx(tx)

	// Create transaction
	transaction := &models.Transaction{
		ID:         uuid.New().String(),
		FromUserID: userID,
		ToUserID:   "",
		Amount:     amount,
		Type:       models.TransactionTypeDeposit,
		Status:     models.TransactionStatusSuccess,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := transactionRepo.Create(transaction); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to create transaction")
	}

	// Update wallet balance
	wallet.Balance += amount
	wallet.UpdatedAt = time.Now()
	if err := walletRepo.Update(wallet); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to update wallet")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to commit transaction")
	}

	s.Cache.Delete(userID)
	return wallet.Balance, nil
}

func (s *walletService) Withdraw(userID string, amount float64) (float64, *APIError) {
	if amount <= 0 {
		return 0, NewBadRequestError("Invalid amount")
	}

	wallet, err := s.WalletRepo.FindByUserID(userID)
	if err != nil {
		return 0, NewInternalServerError("Failed to get wallet")
	}

	if wallet.Balance < amount {
		return 0, NewBadRequestError("Insufficient balance")
	}

	// Start a database transaction for the create and update operations
	tx := s.WalletRepo.DB().Begin()
	if tx.Error != nil {
		return 0, NewInternalServerError("Failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create repository instances with transaction
	walletRepo := s.WalletRepo.WithTx(tx)
	transactionRepo := s.TransactionRepo.WithTx(tx)

	// Create transaction
	transaction := &models.Transaction{
		ID:         uuid.New().String(),
		FromUserID: userID,
		ToUserID:   "",
		Amount:     amount,
		Type:       models.TransactionTypeWithdraw,
		Status:     models.TransactionStatusSuccess,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := transactionRepo.Create(transaction); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to create transaction")
	}

	// Update wallet balance
	wallet.Balance -= amount
	wallet.UpdatedAt = time.Now()
	if err := walletRepo.Update(wallet); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to update wallet")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to commit transaction")
	}

	s.Cache.Delete(userID)
	return wallet.Balance, nil
}

func (s *walletService) Transfer(fromUserID, toUserID string, amount float64) (float64, *APIError) {
	if amount <= 0 {
		return 0, NewBadRequestError("Invalid amount")
	}

	// Get sender's wallet
	fromWallet, err := s.WalletRepo.FindByUserID(fromUserID)
	if err != nil {
		return 0, NewInternalServerError("Failed to get sender's wallet")
	}

	// Get recipient's wallet
	toWallet, err := s.WalletRepo.FindByUserID(toUserID)
	if err != nil {
		return 0, NewInternalServerError("Failed to get recipient's wallet")
	}

	if fromWallet.Balance < amount {
		return 0, NewBadRequestError("Insufficient balance")
	}

	// Start a database transaction
	tx := s.WalletRepo.DB().Begin()
	if tx.Error != nil {
		return 0, NewInternalServerError("Failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create repository instances with transaction
	walletRepo := s.WalletRepo.WithTx(tx)
	transactionRepo := s.TransactionRepo.WithTx(tx)

	// Create transaction
	transaction := &models.Transaction{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Amount:     amount,
		Type:       models.TransactionTypeTransfer,
		Status:     models.TransactionStatusSuccess,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := transactionRepo.Create(transaction); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to create transaction")
	}

	// Update sender's wallet
	fromWallet.Balance -= amount
	fromWallet.UpdatedAt = time.Now()
	if err := walletRepo.Update(fromWallet); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to update sender's wallet")
	}

	// Update recipient's wallet
	toWallet.Balance += amount
	toWallet.UpdatedAt = time.Now()
	if err := walletRepo.Update(toWallet); err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to update recipient's wallet")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, NewInternalServerError("Failed to commit transaction")
	}

	s.Cache.Delete(fromUserID)
	s.Cache.Delete(toUserID)
	return fromWallet.Balance, nil
}

func (s *walletService) GetBalance(userID string) (float64, *APIError) {
	wallet, err := s.WalletRepo.FindByUserID(userID)
	if err != nil {
		return 0, NewInternalServerError("Failed to get wallet")
	}
	return wallet.Balance, nil
}

func (s *walletService) GetTransactionHistory(userID string, page, pageSize int, transactionType, status string) ([]models.Transaction, *APIError) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// Use cache for page=1 and pageSize=10
	if page == 1 && pageSize == 10 && transactionType == "" && status == "" {
		if cachedTransactions, found := s.Cache.Get(userID); found {
			return cachedTransactions.([]models.Transaction), nil
		}
	}

	transactions, err := s.TransactionRepo.FindByUserID(userID, page, pageSize, transactionType, status)
	if err != nil {
		return nil, NewInternalServerError("Failed to get transaction history")
	}

	// Cache the results for page=1 and pageSize=10
	if page == 1 && pageSize == 10 && transactionType == "" && status == "" {
		s.Cache.Set(userID, transactions, 60*time.Minute)
	}

	return transactions, nil
}
