package middleware

import (
	"net/http"
	"time"

	"wallet/internal/repositories"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	UserTokenRepo repositories.UserTokenRepository
	UserRepo      repositories.UserRepository
}

func NewAuthMiddleware(userTokenRepo repositories.UserTokenRepository, userRepo repositories.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		UserTokenRepo: userTokenRepo,
		UserRepo:      userRepo,
	}
}

func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			return
		}

		// Remove "Bearer " prefix if present
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		userToken, err := m.UserTokenRepo.FindByToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if userToken.ExpiresAt.Before(time.Now()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		user, err := m.UserRepo.FindByID(userToken.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
