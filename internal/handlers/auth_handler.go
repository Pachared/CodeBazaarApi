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
	httpx.Fail(c, httpx.NewAppError(http.StatusGone, "route นี้ถูกปิดแล้ว กรุณาใช้ Google sign-in ฝั่ง client แล้วส่ง access token มาที่ /auth/google/session"))
}

func (h *AuthHandler) ExchangeGoogleSession(c *gin.Context) {
	var request contracts.GoogleSessionExchangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Fail(c, httpx.NewAppError(http.StatusBadRequest, "ข้อมูล Google session ไม่ถูกต้อง"))
		return
	}

	response, err := h.authService.ExchangeGoogleSession(request.AccessToken, request.Intent)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
