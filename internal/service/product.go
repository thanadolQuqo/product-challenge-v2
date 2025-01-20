package services

import (
	"context"
	"product-challenge/internal/models"
	"product-challenge/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, req *models.UpsertProductRequest) (*models.Products, error)
	GetAll() ([]models.Products, error)
	GetById(ctx context.Context, id int) (*models.Products, error)
	GetByName(name string) ([]models.Products, error)
	Update(ctx context.Context, product *models.UpsertProductRequest, productId int) (*models.Products, error)
	UpdateStock(ctx context.Context, product *models.ProductStockUpdateReq, productId int) (*models.Products, error)
	Delete(ctx context.Context, id int) error
	DeleteImage(ctx context.Context, id int) error
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

func (s *productService) GetById(ctx context.Context, id int) (*models.Products, error) {
	product, err := s.repo.GetById(ctx, id)
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

func (s *productService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *productService) DeleteImage(ctx context.Context, id int) error {
	err := s.repo.DeleteImage(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *productService) UpdateStock(ctx context.Context, req *models.ProductStockUpdateReq, productId int) (*models.Products, error) {
	res, err := s.repo.UpdateStock(ctx, req, productId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
