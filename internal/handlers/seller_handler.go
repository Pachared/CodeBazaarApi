package handlers

import (
	"net/http"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/middleware"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	sellerService *services.SellerService
}

func NewSellerHandler(sellerService *services.SellerService) *SellerHandler {
	return &SellerHandler{sellerService: sellerService}
}

func (h *SellerHandler) OpenSellerAccount(c *gin.Context) {
	response, err := h.sellerService.OpenSellerAccount(middleware.GetCurrentUser(c))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SellerHandler) SubmitListing(c *gin.Context) {
	var request contracts.SellerListingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Fail(c, httpx.NewAppError(http.StatusBadRequest, "ข้อมูลรายการขายยังไม่ครบหรือไม่ถูกต้อง"))
		return
	}

	response, err := h.sellerService.SubmitListing(middleware.GetCurrentUser(c), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SellerHandler) ListSellerOrders(c *gin.Context) {
	response, err := h.sellerService.ListSellerOrders(middleware.GetCurrentUser(c))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
