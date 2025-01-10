package router

import (
	"github.com/gin-gonic/gin"
	"product-challenge/internal/controller"
	"product-challenge/internal/router/routes"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		engine: gin.Default(),
	}
}

func (r *Router) SetupRoutes(controller *controller.ProductController) {
	// API version group
	v1 := r.engine.Group("/api/v1")

	routes.SetupProductRoutes(v1, *controller)
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
