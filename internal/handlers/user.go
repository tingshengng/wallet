package handlers

import (
	"errors"
	"net/http"
	"time"

	"wallet/internal/models"
	"wallet/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo      repositories.UserRepository
	UserTokenRepo repositories.UserTokenRepository
	WalletRepo    repositories.WalletRepository
}

type LoginRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func NewUserHandler(
	userRepo repositories.UserRepository,
	userTokenRepo repositories.UserTokenRepository,
	walletRepo repositories.WalletRepository,
) *UserHandler {
	return &UserHandler{
		UserRepo:      userRepo,
		UserTokenRepo: userTokenRepo,
		WalletRepo:    walletRepo,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	var userToken *models.UserToken

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if req.Email == "" || req.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Email and Name are required fields"})
		return
	}

	// Find user by email
	user, err := h.UserRepo.FindByEmail(req.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
			return
		}
		// Create user if not exists
		user = &models.User{
			ID:        uuid.New().String(),
			Email:     req.Email,
			Name:      req.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.UserRepo.Create(user); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Create wallet for new user
		wallet := &models.Wallet{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			Balance:   0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.WalletRepo.Create(wallet); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
			return
		}

		// Create new token for new user
		userToken = &models.UserToken{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			Token:     h.generateToken(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(4 * time.Hour),
		}
		if err := h.UserTokenRepo.Create(userToken); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
			return
		}
		c.JSON(http.StatusOK, LoginResponse{
			Token: userToken.Token,
		})
		return
	}

	// Find existing token by user ID
	userToken, err = h.UserTokenRepo.FindByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user token"})
		return
	}

	// Refresh token expiry
	userToken.Token = h.generateToken()
	userToken.ExpiresAt = time.Now().Add(4 * time.Hour)
	userToken.UpdatedAt = time.Now()
	if err := h.UserTokenRepo.Update(userToken); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: userToken.Token,
	})
}

func (h *UserHandler) generateToken() string {
	return "token-" + uuid.New().String()
}
