package main

import (
	"log"

	"github.com/Pachared/CodeBazaarApi/internal/config"
	"github.com/Pachared/CodeBazaarApi/internal/database"
	"github.com/Pachared/CodeBazaarApi/internal/handlers"
	"github.com/Pachared/CodeBazaarApi/internal/repositories"
	"github.com/Pachared/CodeBazaarApi/internal/routes"
	"github.com/Pachared/CodeBazaarApi/internal/services"
	"github.com/Pachared/CodeBazaarApi/internal/session"
)

func main() {
	cfg := config.Load()

	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	if cfg.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("database migration failed: %v", err)
		}
	}

	sessionManager, err := session.NewManager(cfg.SessionSecret, cfg.SessionTTL)
	if err != nil {
		log.Fatalf("session manager setup failed: %v", err)
	}

	userRepository := repositories.NewUserRepository(db)
	productRepository := repositories.NewProductRepository(db)
	orderRepository := repositories.NewOrderRepository(db)
	cookieConsentRepository := repositories.NewCookieConsentRepository(db)

	authService := services.NewAuthService(userRepository, sessionManager)
	catalogService := services.NewCatalogService(productRepository)
	checkoutService := services.NewCheckoutService(db, userRepository, productRepository, orderRepository)
	sellerService := services.NewSellerService(userRepository, productRepository, orderRepository)
	profileService := services.NewProfileService(userRepository)
	downloadsService := services.NewDownloadsService(userRepository, orderRepository)
	cookieService := services.NewCookieService(cookieConsentRepository)

	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(catalogService)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutService)
	sellerHandler := handlers.NewSellerHandler(sellerService)
	profileHandler := handlers.NewProfileHandler(profileService)
	downloadHandler := handlers.NewDownloadHandler(downloadsService)
	cookieHandler := handlers.NewCookieHandler(cookieService)
	healthHandler := handlers.NewHealthHandler()

	router := routes.New(
		cfg,
		userRepository,
		sessionManager,
		authHandler,
		productHandler,
		checkoutHandler,
		sellerHandler,
		profileHandler,
		downloadHandler,
		cookieHandler,
		healthHandler,
	)

	log.Printf("CodeBazaar API listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
