package routes

import (
	"net/http"

	"github.com/Pachared/CodeBazaarApi/internal/config"
	"github.com/Pachared/CodeBazaarApi/internal/handlers"
	"github.com/Pachared/CodeBazaarApi/internal/middleware"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
	"github.com/Pachared/CodeBazaarApi/internal/session"
	"github.com/gin-gonic/gin"
)

func New(
	cfg config.Config,
	userRepository *repositories.UserRepository,
	sessionManager *session.Manager,
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	checkoutHandler *handlers.CheckoutHandler,
	sellerHandler *handlers.SellerHandler,
	profileHandler *handlers.ProfileHandler,
	downloadHandler *handlers.DownloadHandler,
	cookieHandler *handlers.CookieHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS(cfg.AllowedOrigins))
	router.Use(middleware.CurrentUser(userRepository, sessionManager))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "CodeBazaar API",
			"version": "v1",
		})
	})

	router.GET("/health", healthHandler.Health)

	mountAPI := func(group *gin.RouterGroup) {
		group.POST("/auth/google/start", authHandler.StartGoogleAuth)
		group.POST("/auth/google/session", authHandler.ExchangeGoogleSession)

		group.GET("/products", productHandler.ListProducts)
		group.GET("/products/featured", productHandler.ListFeaturedProducts)
		group.GET("/products/:productID", productHandler.GetProductByID)

		group.GET("/sellers", productHandler.ListSellers)
		group.GET("/sellers/:sellerSlug", productHandler.GetSellerBySlug)
		group.GET("/sellers/:sellerSlug/products", productHandler.ListSellerProducts)

		group.POST("/checkout/orders", checkoutHandler.SubmitOrder)

		group.POST("/seller/onboarding/google", sellerHandler.OpenSellerAccount)
		group.POST("/seller/onboarding/github", sellerHandler.OpenSellerAccount)
		group.POST("/seller/listings", sellerHandler.SubmitListing)
		group.GET("/seller/orders", sellerHandler.ListSellerOrders)

		group.GET("/me/profile", profileHandler.GetProfile)
		group.PUT("/me/profile", profileHandler.UpdateProfile)

		group.GET("/me/downloads", downloadHandler.ListDownloads)
		group.POST("/me/downloads/:libraryItemID/download", downloadHandler.MarkDownloaded)

		group.GET("/cookie-consent", cookieHandler.GetConsent)
		group.PUT("/cookie-consent", cookieHandler.SaveConsent)
		group.GET("/me/cookie-consent", cookieHandler.GetConsent)
		group.PUT("/me/cookie-consent", cookieHandler.SaveConsent)
	}

	mountAPI(router.Group(""))
	mountAPI(router.Group("/api/v1"))

	return router
}
