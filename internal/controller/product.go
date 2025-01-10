package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"product-challenge/internal/models"
	"product-challenge/internal/service"
	"strconv"
	"time"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(service services.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var product models.UpsertProductRequest
	if err := ctx.ShouldBind(&product); err != nil {
		fmt.Println("error should bind")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the uploaded file from the form-data
	file, fileHeader, err := ctx.Request.FormFile("image")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		fmt.Println("error get image : ", err)

	}

	if fileHeader != nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		product.Filename = fileName
		product.ImageFile = &file
	}

	createdProduct, err := c.service.Create(ctx, &product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdProduct)
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	products, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (c *ProductController) GetById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	product, err := c.service.GetById(int(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) GetByName(ctx *gin.Context) {
	name := ctx.Query("name")

	if name == "" {
		ctx.JSON(http.StatusOK, []models.Products{})
		return
	}

	product, err := c.service.GetByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	var product models.UpsertProductRequest
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctx.ShouldBind(&product); err != nil {
		fmt.Println("error should bind")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the uploaded file from the form-data
	file, fileHeader, err := ctx.Request.FormFile("image")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	if fileHeader != nil {
		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
		product.Filename = fileName
		product.ImageFile = &file
	}

	updateProduct, err := c.service.Update(ctx, &product, int(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updateProduct)
}

// DeleteProduct
func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.service.Delete(int(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Product with id %d deleted", id)})
}

func (c *ProductController) DeleteProductImage(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.service.DeleteImage(int(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Image of product Id %d is deleted.", id)})

}
