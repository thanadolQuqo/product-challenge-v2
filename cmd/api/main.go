package main

import (
	"context"
	"log"
	"product-challenge/internal/controller"
	"product-challenge/internal/repository"
	productRouter "product-challenge/internal/router"
	"product-challenge/internal/service"
	"product-challenge/pkg/config"
	"product-challenge/pkg/database"
	"product-challenge/pkg/storage"
)

func main() {
	// load configs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// connect to DB
	db, err := database.NewCockroachDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ctx := context.Background()
	aws, err := storage.NewS3Client(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to get aws client: %v", err)
	}

	productRepo := repository.NewProductRepository(db, aws, cfg)
	productService := services.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	// initialize router
	r := productRouter.NewRouter()

	// setup router
	r.SetupRoutes(productController)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
