package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"product-challenge/internal/models"
	"product-challenge/pkg/config"
	"time"
)

type ProductRepository interface {
	Create(ctx context.Context, products *models.UpsertProductRequest) (*models.Products, error)
	GetAll() ([]models.Products, error)
	GetById(ctx context.Context, productId int) (*models.Products, error)
	GetByName(name string) ([]models.Products, error)
	Update(ctx context.Context, productReq *models.UpsertProductRequest, productId int) (*models.Products, error)
	Delete(ctx context.Context, productId int) error
	DeleteImage(ctx context.Context, productId int) error
}

type productRepository struct {
	db        *gorm.DB
	awsClient *s3.Client
	cfg       config.Config
	redis     *redis.Client
}

func NewProductRepository(db *gorm.DB, client *s3.Client, cfg *config.Config, redis *redis.Client) ProductRepository {
	return &productRepository{db: db, awsClient: client, cfg: *cfg, redis: redis}
}

func (r *productRepository) Create(ctx context.Context, req *models.UpsertProductRequest) (*models.Products, error) {
	productData := models.Products{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
	}

	if req.ImageFile != nil {
		// 1. upload file
		imageUrl, err := r.UploadImage(ctx, req.Filename, *req.ImageFile)
		if err != nil {
			return nil, err
		}

		// set image url to product data for insert
		productData.ImageName = req.Filename
		productData.ImageURL = imageUrl
	}

	// 2. insert into db
	if err := r.db.Create(&productData).Error; err != nil {
		return nil, err
	}

	return &productData, nil
}

func (r *productRepository) GetAll() ([]models.Products, error) {
	// DISCUSSION: should we cache this one?

	var products []models.Products
	err := r.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetById(ctx context.Context, productId int) (*models.Products, error) {
	var product models.Products

	// 1. if data from redis available, use it
	key := fmt.Sprintf("product:id:%d", productId)
	productCache, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			fmt.Println("get data from redis error: ", err)
			return nil, err
		}
	} else { // if found, return
		err = json.Unmarshal([]byte(productCache), &product)
		if err != nil {
			return nil, err
		}
		return &product, nil
	}

	// 2. if not, query
	err = r.db.First(&product, productId).Error
	if err != nil {
		return nil, err
	}

	// Cache the product to Redis
	productJSON, _ := json.Marshal(product)
	if err = r.redis.Set(ctx, key, productJSON, time.Hour*10).Err(); err != nil {
		fmt.Println("error save to redis : ", err)
	}
	return &product, nil
}

func (r *productRepository) GetByName(name string) ([]models.Products, error) {
	var products []models.Products
	err := r.db.Where("name ILIKE ?", name+"%").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(ctx context.Context, req *models.UpsertProductRequest, productId int) (*models.Products, error) {
	// 1. check if there is already existing image
	var product models.Products
	err := r.db.First(&product, productId).Error
	if err != nil {
		return nil, err
	}

	// if image name and image url exist
	if product.ImageName != "" && product.ImageURL != "" {
		err = r.DeleteImage(ctx, productId)
		if err != nil {
			return nil, err
		}
	}

	// set data for update
	product = models.Products{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Stock:       req.Stock,
	}

	if req.ImageFile != nil {
		// 1. upload file
		imageUrl, err := r.UploadImage(ctx, req.Filename, *req.ImageFile)
		if err != nil {
			return nil, err
		}
		// set image url to product data for insert
		product.ImageName = req.Filename
		product.ImageURL = imageUrl
	}

	// 2. insert into db
	if err := r.db.Where("id = ?", productId).
		Updates(&product).Error; err != nil {
		return nil, err
	}

	err = r.RemoveProductCache(ctx, productId)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Delete(ctx context.Context, productId int) error {
	if err := r.db.Delete(&models.Products{}, "id = ?", productId).Error; err != nil {
		return err
	}

	err := r.RemoveProductCache(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) DeleteImage(ctx context.Context, productId int) error {
	// 1. get S3 link
	var product models.Products
	err := r.db.First(&product, productId).Error
	if err != nil {
		return err
	}

	// 2. delete from S3 bucket
	// Create the input for the DeleteObject API
	input := &s3.DeleteObjectInput{
		Bucket: &r.cfg.Aws.BucketName, // The S3 bucket name
		Key:    &product.ImageName,    // The key of the object to delete (object name)
	}

	_, err = r.awsClient.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	// 3. clear image record
	// clear image name and URL
	product.ImageName = ""
	product.ImageURL = ""

	if err := r.db.Where("id = ?", productId).
		Save(&product).Error; err != nil {
		return err
	}

	err = r.RemoveProductCache(ctx, productId)
	if err != nil {
		return err
	}
	return nil
}

func (r *productRepository) UploadImage(ctx context.Context, fileName string, file multipart.File) (string, error) {
	client := r.awsClient
	region := r.cfg.Aws.Region
	bucket := r.cfg.Aws.BucketName
	// Determine the content type
	contentType := http.DetectContentType(make([]byte, 512))
	_, err := file.Read(make([]byte, 512)) // Pre-read a portion of the file
	if err != nil {
		contentType = "application/octet-stream"
	}

	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %v", err)
	}

	// Upload the file
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        file, // Convert to a reader
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Generate the S3 public URL
	s3URL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, fileName)

	return s3URL, nil
}

// RemoveProductCache is a util function that help remove product from cache by feeding product id to it.
func (r *productRepository) RemoveProductCache(ctx context.Context, productId int) error {
	// 1. check if that product Id is in cache. if it is, clear from cache
	key := fmt.Sprintf("product:id:%d", productId)
	_, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			fmt.Println("get data from redis error: ", err)
			return err
		}
		fmt.Println("cache of this product is not available")
	} else { // if found, clear that value from redis
		err = r.redis.Del(ctx, key).Err()
		if err != nil {
			fmt.Println("delete data from redis error : ", err)
		}
	}
	return nil
}
