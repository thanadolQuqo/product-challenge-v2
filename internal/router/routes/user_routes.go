package routes

import (
	"github.com/gin-gonic/gin"
	"product-challenge/internal/controller"
)

func SetupUserRoutes(router *gin.RouterGroup, controller controller.UserController) {
	products := router.Group("/user")
	{
		// create
		products.POST("/register", controller.UserRegister)
		products.POST("/login", controller.UserLogin)
	}
}
