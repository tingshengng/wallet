package repositories

import (
	"wallet/internal/models"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByUserID(userID string, page, pageSize int, transactionType, status string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID)

	if transactionType != "" {
		query = query.Where("type = ?", transactionType)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id string) error {
	return r.db.Delete(&models.Transaction{}, "id = ?", id).Error
}

func (r *transactionRepository) WithTx(tx interface{}) TransactionRepository {
	txDB, ok := tx.(*gorm.DB)
	if !ok {
		return r
	}
	return &transactionRepository{db: txDB}
}
