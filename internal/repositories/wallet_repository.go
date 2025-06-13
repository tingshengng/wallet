package repositories

import (
	"wallet/internal/models"

	"gorm.io/gorm"
)

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) FindByUserID(userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *walletRepository) Delete(id string) error {
	return r.db.Delete(&models.Wallet{}, "id = ?", id).Error
}

func (r *walletRepository) WithTx(tx interface{}) WalletRepository {
	txDB, ok := tx.(*gorm.DB)
	if !ok {
		return r
	}
	return &walletRepository{db: txDB}
}

// DB returns the underlying GORM DB instance
func (r *walletRepository) DB() *gorm.DB {
	return r.db
}
