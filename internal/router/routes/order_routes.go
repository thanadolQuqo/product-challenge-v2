package routes

import (
	"github.com/gin-gonic/gin"
	"product-challenge/internal/controller"
)

func SetupOrderRoutes(router *gin.RouterGroup, middleware gin.HandlerFunc, controller controller.OrderController) {
	carts := router.Group("/carts")
	carts.Use(middleware)
	{
		// get cart info with summary
		carts.GET("", controller.GetCart)

		// add item to cart
		carts.POST("/add", controller.AddProductToCart)

		// remove item from cart
		carts.DELETE("/remove", controller.RemoveProductFromCart)

	}

	orders := router.Group("/orders")
	orders.Use(middleware)
	{
		// make order using all items in cart
		orders.POST("/make", controller.MakeOrder)
	}
}
