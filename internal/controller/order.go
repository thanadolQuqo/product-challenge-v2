package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"product-challenge/internal/models"
	services "product-challenge/internal/service"
)

type OrderController struct {
	service services.OrderService
}

func NewOrderController(service services.OrderService) *OrderController {
	return &OrderController{service: service}
}

func (c *OrderController) GetCart(ctx *gin.Context) {
	name := ctx.Query("username")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	createResp, err := c.service.GetCart(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createResp)
}

func (c *OrderController) AddProductToCart(ctx *gin.Context) {
	var req models.CartRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resp, err := c.service.AddProductToCart(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *OrderController) RemoveProductFromCart(ctx *gin.Context) {
	var req models.CartRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resp, err := c.service.RemoveProductFromCart(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *OrderController) MakeOrder(ctx *gin.Context) {

	name := ctx.Query("username")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}
	loginResp, err := c.service.MakeOrder(ctx, name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loginResp)
}
