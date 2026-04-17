package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, message string) error {
	return &AppError{
		Status:  status,
		Message: message,
	}
}

func Fail(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "เกิดข้อผิดพลาดบางอย่างจากเซิร์ฟเวอร์"

	var appError *AppError
	if errors.As(err, &appError) {
		status = appError.Status
		message = appError.Message
	} else if err != nil && err.Error() != "" {
		message = err.Error()
	}

	c.AbortWithStatusJSON(status, gin.H{
		"message": message,
	})
}
