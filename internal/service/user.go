package services

import (
	"context"
	"product-challenge/internal/models"
	"product-challenge/internal/repository"
)

type UserService interface {
	UserRegister(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error)
	UserLogin(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) UserRegister(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error) {
	resp, err := s.repo.UserRegister(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *userService) UserLogin(ctx context.Context, req *models.UserAuthRequest) (*models.UserAuthResponse, error) {
	resp, err := s.repo.UserLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
