package handlers

import (
	"net/http"

	"github.com/Pachared/CodeBazaarApi/internal/contracts"
	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/middleware"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	checkoutService *services.CheckoutService
}

func NewCheckoutHandler(checkoutService *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutService: checkoutService}
}

func (h *CheckoutHandler) SubmitOrder(c *gin.Context) {
	var request contracts.CheckoutSubmitInput
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Fail(c, httpx.NewAppError(http.StatusBadRequest, "ข้อมูลคำสั่งซื้อไม่ครบหรือไม่ถูกต้อง"))
		return
	}

	response, err := h.checkoutService.SubmitOrder(middleware.GetCurrentUser(c), request)
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
