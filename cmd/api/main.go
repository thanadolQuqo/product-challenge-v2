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
	// connect to redis
	redis, err := database.NewRedisCache(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	aws, err := storage.NewS3Client(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to get aws client: %v", err)
	}

	productRepo := repository.NewProductRepository(db, aws, cfg, redis)
	productService := services.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	userRepo := repository.NewUserRepository(db, cfg)
	userService := services.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	orderRepo := repository.NewOrderRepository(db, cfg)
	orderService := services.NewOrderService(orderRepo)
	orderController := controller.NewOrderController(orderService)
	// initialize router
	r := productRouter.NewRouter(cfg)

	// setup router
	r.SetupRoutes(productController, userController, orderController)

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
