package repository

import (
	"gorm.io/gorm"
	"product-challenge/internal/models"
)

type ProductRepository interface {
	Create(products *models.Products) error
	GetAll() ([]models.Products, error)
	GetById(id int) (*models.Products, error)
	GetByName(name string) ([]models.Products, error)
	Update(product *models.Products) (*models.Products, error)
	Delete(id int) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Products) error {
	return r.db.Create(product).Error
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

func (r *productRepository) Update(product *models.Products) (*models.Products, error) {
	if err := r.db.Where("id = ?", product.ID).
		Updates(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) Delete(id int) error {
	if err := r.db.Delete(&models.Products{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
