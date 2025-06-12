package repositories

import (
	"wallet/internal/models"
	"gorm.io/gorm"
)

type userTokenRepository struct {
	db *gorm.DB
}

func NewUserTokenRepository(db *gorm.DB) UserTokenRepository {
	return &userTokenRepository{db: db}
}

func (r *userTokenRepository) Create(token *models.UserToken) error {
	return r.db.Create(token).Error
}

func (r *userTokenRepository) FindByToken(token string) (*models.UserToken, error) {
	var userToken models.UserToken
	if err := r.db.Where("token = ?", token).First(&userToken).Error; err != nil {
		return nil, err
	}
	return &userToken, nil
}

func (r *userTokenRepository) FindByUserID(userID string) (*models.UserToken, error) {
	var userToken models.UserToken
	if err := r.db.Where("user_id = ?", userID).First(&userToken).Error; err != nil {
		return nil, err
	}
	return &userToken, nil
}

func (r *userTokenRepository) Update(token *models.UserToken) error {
	return r.db.Save(token).Error
}

func (r *userTokenRepository) Delete(id string) error {
	return r.db.Delete(&models.UserToken{}, id).Error
}
