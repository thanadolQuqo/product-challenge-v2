package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"product-challenge/internal/controller"
	"product-challenge/internal/router/routes"
	"product-challenge/pkg/config"
)

type Router struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewRouter(cfg *config.Config) *Router {
	return &Router{
		engine: gin.Default(),
		cfg:    cfg,
	}
}

func (r *Router) SetupRoutes(productController *controller.ProductController, userController *controller.UserController, orderController *controller.OrderController) {
	// API version group
	v1 := r.engine.Group("/api/v1")

	routes.SetupUserRoutes(v1, *userController)

	// adding middleware to check for token only in product api
	routes.SetupProductRoutes(v1, r.authMiddleware(), *productController)

	routes.SetupOrderRoutes(v1, r.authMiddleware(), *orderController)
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func (r *Router) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(r.cfg.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Stop further processing if unauthorized
			return
		}

		// Set the token claims to the context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next() // Proceed to the next handler if authorized
	}
}
