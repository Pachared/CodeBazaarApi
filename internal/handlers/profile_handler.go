package handlers

import (
	"net/http"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/middleware"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *services.ProfileService
}

func NewProfileHandler(profileService *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	response, err := h.profileService.GetProfile(middleware.GetCurrentUser(c))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var request contracts.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Fail(c, httpx.NewAppError(http.StatusBadRequest, "ข้อมูลโปรไฟล์ไม่ถูกต้อง"))
		return
	}

	response, err := h.profileService.UpdateProfile(middleware.GetCurrentUser(c), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
