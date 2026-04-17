package handlers

import (
	"net/http"

	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/middleware"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	downloadsService *services.DownloadsService
}

func NewDownloadHandler(downloadsService *services.DownloadsService) *DownloadHandler {
	return &DownloadHandler{downloadsService: downloadsService}
}

func (h *DownloadHandler) ListDownloads(c *gin.Context) {
	response, err := h.downloadsService.ListDownloads(middleware.GetCurrentUser(c))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *DownloadHandler) MarkDownloaded(c *gin.Context) {
	response, err := h.downloadsService.MarkDownloaded(middleware.GetCurrentUser(c), c.Param("libraryItemID"))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
