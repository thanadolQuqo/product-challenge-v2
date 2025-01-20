package services

import (
	"context"
	"product-challenge/internal/models"
	"product-challenge/internal/repository"
)

type OrderService interface {
	GetCart(ctx context.Context, username string) (*models.GetCartResponse, error)
	AddProductToCart(ctx context.Context, req *models.CartRequest) (bool, error)
	RemoveProductFromCart(ctx context.Context, req *models.CartRequest) (bool, error)
	MakeOrder(ctx context.Context, username string) (*models.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) GetCart(ctx context.Context, username string) (*models.GetCartResponse, error) {
	resp, err := s.repo.GetCart(ctx, username)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *orderService) AddProductToCart(ctx context.Context, req *models.CartRequest) (bool, error) {
	resp, err := s.repo.AddProductToCart(ctx, req)
	if err != nil {
		return false, err
	}

	return resp, nil
}

func (s *orderService) RemoveProductFromCart(ctx context.Context, req *models.CartRequest) (bool, error) {
	resp, err := s.repo.RemoveProductFromCart(ctx, req)
	if err != nil {
		return false, err
	}

	return resp, nil
}

func (s *orderService) MakeOrder(ctx context.Context, username string) (*models.Order, error) {
	resp, err := s.repo.MakeOrder(ctx, username)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
