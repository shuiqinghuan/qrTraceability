package main

import (
	"fmt"
	"log"

	"qr-traceability/internal/config"
	"qr-traceability/internal/handlers"
	"qr-traceability/internal/models"
	"qr-traceability/internal/repository"
	"qr-traceability/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := config.InitDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Auto migrate database models
	if err := db.AutoMigrate(
		&models.Product{},
		&models.ProductBatch{},
		&models.MediaFile{},
		&models.HarvestQuality{},
		&models.UserInteraction{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")

	// Initialize Redis
	redisClient, err := config.InitRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Set up Gin router
	router := gin.Default()

	// Initialize repositories
	productRepo := repository.NewProductRepository(db)
	batchRepo := repository.NewBatchRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	harvestRepo := repository.NewHarvestRepository(db)
	interactionRepo := repository.NewInteractionRepository(db)

	// Initialize services
	productService := service.NewProductService(productRepo)
	batchService := service.NewBatchService(batchRepo, mediaRepo, harvestRepo)
	interactionService := service.NewInteractionService(interactionRepo, redisClient)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	batchHandler := handlers.NewBatchHandler(batchService)
	interactionHandler := handlers.NewInteractionHandler(interactionService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Product routes
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("/:id", productHandler.GetProductByID)
			products.GET("/code/:code", productHandler.GetProductByCode)
			products.GET("", productHandler.ListProducts)
			products.PUT("/:id", productHandler.UpdateProduct)
			products.DELETE("/:id", productHandler.DeleteProduct)
		}

		// Batch routes
		batches := api.Group("/batches")
		{
			batches.POST("", batchHandler.CreateBatch)
			batches.GET("/:id", batchHandler.GetBatchByID)
			batches.GET("/unique/:unique_id", batchHandler.GetBatchByUniqueID)
			batches.GET("/product/:product_id", batchHandler.ListBatchesByProductID)
			batches.PUT("/:id", batchHandler.UpdateBatch)
			batches.DELETE("/:id", batchHandler.DeleteBatch)
			batches.POST("/:batch_id/media", batchHandler.AddMedia)
			batches.POST("/:batch_id/harvest", batchHandler.UpdateHarvestQuality)
			batches.GET("/:batch_id/qr", batchHandler.GenerateQRCode)
		}

		// Interaction routes
		interactions := api.Group("/interactions")
		{
			interactions.POST("", interactionHandler.RecordInteraction)
			interactions.GET("/batch/:batch_id", interactionHandler.GetInteractionStats)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}