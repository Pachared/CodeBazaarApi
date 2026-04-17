package handlers

import (
	"net/http"
	"strings"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/middleware"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type CookieHandler struct {
	cookieService *services.CookieService
}

func NewCookieHandler(cookieService *services.CookieService) *CookieHandler {
	return &CookieHandler{cookieService: cookieService}
}

func (h *CookieHandler) GetConsent(c *gin.Context) {
	response, err := h.cookieService.GetConsent(
		middleware.GetCurrentUser(c),
		strings.TrimSpace(c.GetHeader("X-Session-Key")),
	)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CookieHandler) SaveConsent(c *gin.Context) {
	var request contracts.CookieConsentUpsertRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Fail(c, httpx.NewAppError(http.StatusBadRequest, "ข้อมูล cookie consent ไม่ถูกต้อง"))
		return
	}

	response, err := h.cookieService.SaveConsent(
		middleware.GetCurrentUser(c),
		strings.TrimSpace(c.GetHeader("X-Session-Key")),
		request,
	)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
