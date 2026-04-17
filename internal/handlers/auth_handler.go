package handlers

import (
	"net/http"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) StartGoogleAuth(c *gin.Context) {
	var request contracts.AuthStartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Fail(c, httpx.NewAppError(http.StatusBadRequest, "รูปแบบข้อมูลเข้าสู่ระบบไม่ถูกต้อง"))
		return
	}

	response, err := h.authService.StartGoogleAuth(request.Intent)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
