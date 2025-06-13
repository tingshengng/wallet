package repositories

import (
	"wallet/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
}

type UserTokenRepository interface {
	Create(token *models.UserToken) error
	FindByToken(token string) (*models.UserToken, error)
	FindByUserID(userID string) (*models.UserToken, error)
	Update(token *models.UserToken) error
	Delete(id string) error
}

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByUserID(userID string) (*models.Wallet, error)
	Update(wallet *models.Wallet) error
	Delete(id string) error
	WithTx(tx interface{}) WalletRepository
	DB() *gorm.DB
}

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByUserID(userID string, page, pageSize int, transactionType, status string) ([]models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id string) error
	WithTx(tx interface{}) TransactionRepository
}
