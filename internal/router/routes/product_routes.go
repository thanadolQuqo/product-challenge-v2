package routes

import (
	"github.com/gin-gonic/gin"
	"product-challenge/internal/controller"
)

func SetupProductRoutes(router *gin.RouterGroup, middleware gin.HandlerFunc, controller controller.ProductController) {
	products := router.Group("/products")
	products.Use(middleware)
	{
		// read
		products.GET("", controller.GetAllProducts)
		products.GET("/:id", controller.GetById)
		products.GET("/search", controller.GetByName)
		// create
		products.POST("", controller.CreateProduct)
		// update
		products.PUT("/:id", controller.UpdateProduct)
		//delete
		products.DELETE("/:id", controller.DeleteProduct)
		products.DELETE("/image/:id", controller.DeleteProductImage)
	}
}
