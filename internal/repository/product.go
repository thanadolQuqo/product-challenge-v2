package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/gorm"
	"mime/multipart"
	"net/http"
	"product-challenge/internal/models"
	"product-challenge/pkg/config"
)

type ProductRepository interface {
	Create(ctx context.Context, products *models.UpsertProductRequest) (*models.Products, error)
	GetAll() ([]models.Products, error)
	GetById(id int) (*models.Products, error)
	GetByName(name string) ([]models.Products, error)
	Update(ctx context.Context, productReq *models.UpsertProductRequest, productId int) (*models.Products, error)
	Delete(id int) error
	DeleteImage(id int) error
}

type productRepository struct {
	db        *gorm.DB
	awsClient *s3.Client
	cfg       config.Config
}

func NewProductRepository(db *gorm.DB, client *s3.Client, cfg *config.Config) ProductRepository {
	return &productRepository{db: db, awsClient: client, cfg: *cfg}
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
	var users []models.Products
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *productRepository) GetById(id int) (*models.Products, error) {
	var product models.Products
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
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
		err = r.DeleteImage(productId)
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
	//if err := r.db.Create(&productData).Error; err != nil {
	if err := r.db.Where("id = ?", productId).
		Updates(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Delete(id int) error {
	if err := r.db.Delete(&models.Products{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *productRepository) DeleteImage(id int) error {
	// 1. get S3 link
	var product models.Products
	err := r.db.First(&product, id).Error
	if err != nil {
		return err
	}

	// 2. delete from S3 bucket
	// Create the input for the DeleteObject API
	input := &s3.DeleteObjectInput{
		Bucket: &r.cfg.Aws.BucketName, // The S3 bucket name
		Key:    &product.ImageName,    // The key of the object to delete (object name)
	}

	_, err = r.awsClient.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	// 3. clear image record
	// clear image name and URL
	product.ImageName = ""
	product.ImageURL = ""

	if err := r.db.Where("id = ?", id).
		Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) UploadImage(ctx context.Context, fileName string, file multipart.File) (string, error) {
	//region := os.Getenv("AWS_REGION")
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
