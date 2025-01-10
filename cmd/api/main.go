package main

import (
	"log"
	"product-challenge/internal/controller"
	"product-challenge/internal/repository"
	productRouter "product-challenge/internal/router"
	"product-challenge/internal/service"
	"product-challenge/pkg/database"
)

func main() {
	// connect to DB
	db, err := database.NewCockroachDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	productRepo := repository.NewProductRepository(db)
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
