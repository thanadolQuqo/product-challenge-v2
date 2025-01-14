package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"product-challenge/internal/models"
	services "product-challenge/internal/service"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) UserRegister(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if username == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"registration error : ": fmt.Errorf("username or password is empty")})
	}

	// TODO: hash password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	regisReq := models.UserAuthRequest{
		Username: username,
		Password: string(hashedPassword),
	}

	createResp, err := c.service.UserRegister(ctx, &regisReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createResp)
}

func (c *UserController) UserLogin(ctx *gin.Context) {

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")

	if username == "" || password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"login error : ": fmt.Errorf("username or password is empty")})
	}

	loginReq := models.UserAuthRequest{
		Username: username,
		Password: password,
	}

	loginResp, err := c.service.UserLogin(ctx, &loginReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, loginResp)
}
