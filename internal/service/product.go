package services

import (
	"context"
	"product-challenge/internal/models"
	"product-challenge/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, user *models.UpsertProductRequest) (*models.Products, error)
	GetAll() ([]models.Products, error)
	GetById(id int) (*models.Products, error)
	GetByName(name string) ([]models.Products, error)
	Update(ctx context.Context, product *models.UpsertProductRequest, productId int) (*models.Products, error)
	Delete(id int) error
	DeleteImage(id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, req *models.UpsertProductRequest) (*models.Products, error) {
	product, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) GetAll() ([]models.Products, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) GetById(id int) (*models.Products, error) {
	product, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetByName(name string) ([]models.Products, error) {
	products, err := s.repo.GetByName(name)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) Update(ctx context.Context, req *models.UpsertProductRequest, productId int) (*models.Products, error) {
	product, err := s.repo.Update(ctx, req, productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Delete(id int) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *productService) DeleteImage(id int) error {
	err := s.repo.DeleteImage(id)
	if err != nil {
		return err
	}
	return nil
}
