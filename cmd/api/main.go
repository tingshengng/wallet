package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"wallet/internal/cache"
	"wallet/internal/database"
	"wallet/internal/handlers"
	"wallet/internal/middleware"
	"wallet/internal/migrations"
	"wallet/internal/repositories"
	"wallet/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	// Run migrations
	migrator := migrations.NewMigrator(db)
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}

	userRepo := repositories.NewUserRepository(db)
	userTokenRepo := repositories.NewUserTokenRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	cache := cache.NewInMemoryCache()

	service := services.NewWalletService(walletRepo, transactionRepo, cache)

	userHandler := handlers.NewUserHandler(userRepo, userTokenRepo, walletRepo)
	walletHandler := handlers.NewWalletHandler(service)
	authMiddleware := middleware.NewAuthMiddleware(userTokenRepo, userRepo)

	r := gin.Default()

	// Public routes
	public := r.Group("/api")
	public.POST("/login", userHandler.Login)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(authMiddleware.AuthMiddleware())
	{
		protected.POST("/deposit", walletHandler.Deposit)
		protected.POST("/withdraw", walletHandler.Withdraw)
		protected.POST("/transfer", walletHandler.Transfer)
		protected.GET("/balance", walletHandler.GetBalance)
		protected.GET("/transactions", walletHandler.GetTransactionHistory)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
