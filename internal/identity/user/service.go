package user

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Service struct {
	userRepository Repository
}

func NewService(userRepository Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) *shared_errors.AppError {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	user := &UserEntity{
		Email:    email,
		Password: string(hash),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context, page int) ([]UserEntity, *dtos.PaginationDTO, *shared_errors.AppError) {
	users, pagination, err := s.userRepository.List(ctx, page)
	if err != nil {
		return nil, nil, shared_errors.NewBadRequest(err.Error())
	}
	return users, pagination, nil
}
