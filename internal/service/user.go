package service

import (
	"context"
	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	repository *repository.Queries
}

func NewUserService(repo *repository.Queries) *UserService {
	return &UserService{repository: repo}
}

func (s *UserService) FindById(ctx context.Context, id pgtype.UUID) (*repository.User, error) {
	user, err := s.repository.UserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, params repository.UserCreateParams) (string, error) {
	uuid, err := s.repository.UserCreate(ctx, params)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
