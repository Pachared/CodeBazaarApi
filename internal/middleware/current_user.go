package middleware

import (
	"errors"
	"strings"

	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const currentUserKey = "currentUser"

func CurrentUser(userRepository *repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := strings.TrimSpace(c.GetHeader("X-User-ID"))
		email := strings.TrimSpace(c.GetHeader("X-User-Email"))

		if userID == "" && email == "" {
			c.Next()
			return
		}

		user, err := userRepository.GetByIDOrEmail(userID, email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.Next()
			return
		}

		if user != nil {
			c.Set(currentUserKey, user)
		}

		c.Next()
	}
}

func GetCurrentUser(c *gin.Context) *models.User {
	value, exists := c.Get(currentUserKey)
	if !exists {
		return nil
	}

	user, ok := value.(*models.User)
	if !ok {
		return nil
	}

	return user
}
