package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/models"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
	"github.com/Pachared/CodeBazaarApi/internal/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const currentUserKey = "currentUser"

func CurrentUser(userRepository *repositories.UserRepository, sessionManager *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := sessionTokenFromRequest(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := sessionManager.Parse(token)
		if err != nil {
			httpx.Fail(c, httpx.NewAppError(http.StatusUnauthorized, "session ไม่ถูกต้องหรือหมดอายุแล้ว กรุณาเข้าสู่ระบบใหม่"))
			return
		}

		user, err := userRepository.GetByIDOrEmail(claims.UserID, claims.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.Fail(c, err)
			return
		}

		if user == nil {
			user, err = userRepository.FindOrCreateExternalUser(
				claims.UserID,
				claims.Email,
				claims.Name,
				claims.Provider,
				claims.Role,
			)
			if err != nil {
				httpx.Fail(c, err)
				return
			}
		}

		c.Set(currentUserKey, user)
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

func sessionTokenFromRequest(c *gin.Context) string {
	authorizationHeader := strings.TrimSpace(c.GetHeader("Authorization"))
	if authorizationHeader != "" {
		parts := strings.Fields(authorizationHeader)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return strings.TrimSpace(parts[1])
		}
	}

	return strings.TrimSpace(c.GetHeader("X-Session-Token"))
}
