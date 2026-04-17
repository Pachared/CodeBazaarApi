package handlers

import (
	"net/http"
	"strings"

	"github.com/Pachared/CodeBazaarApi/internal/httpx"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	catalogService *services.CatalogService
}

func NewProductHandler(catalogService *services.CatalogService) *ProductHandler {
	return &ProductHandler{catalogService: catalogService}
}

func (h *ProductHandler) ListFeaturedProducts(c *gin.Context) {
	response, err := h.catalogService.ListFeaturedProducts()
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	response, err := h.catalogService.GetProductByID(c.Param("productID"))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	stacks := c.QueryArray("stack")
	if rawStacks := strings.TrimSpace(c.Query("stacks")); rawStacks != "" {
		stacks = append(stacks, strings.Split(rawStacks, ",")...)
	}

	for index, stack := range stacks {
		stacks[index] = strings.TrimSpace(stack)
	}

	response, err := h.catalogService.ListProducts(services.ProductFilter{
		Query:        c.Query("query"),
		Category:     c.DefaultQuery("category", "all"),
		License:      c.DefaultQuery("license", "all"),
		Price:        c.DefaultQuery("price", "all"),
		Sort:         c.DefaultQuery("sort", "featured"),
		VerifiedOnly: c.Query("verifiedOnly") == "true",
		Stacks:       stacks,
	})
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) ListSellers(c *gin.Context) {
	response, err := h.catalogService.ListSellers()
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetSellerBySlug(c *gin.Context) {
	response, err := h.catalogService.GetSellerBySlug(c.Param("sellerSlug"))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) ListSellerProducts(c *gin.Context) {
	response, err := h.catalogService.ListProductsBySellerSlug(c.Param("sellerSlug"))
	if err != nil {
		httpx.Fail(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
