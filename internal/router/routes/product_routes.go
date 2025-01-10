package routes

import (
	"github.com/gin-gonic/gin"
	"product-challenge/internal/controller"
)

func SetupProductRoutes(router *gin.RouterGroup, controller controller.ProductController) {
	products := router.Group("/products")
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
	}
}
