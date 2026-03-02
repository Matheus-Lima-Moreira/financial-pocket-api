package user

import (
	"context"

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

func (s *Service) List(ctx context.Context, page int) ([]UserEntity, *dtos.PaginationDTO, *shared_errors.AppError) {
	users, pagination, err := s.userRepository.List(ctx, page)
	if err != nil {
		return nil, nil, shared_errors.NewBadRequest(err.Error())
	}
	return users, pagination, nil
}

func (s *Service) Details(ctx context.Context, id string) (*UserEntity, *shared_errors.AppError) {
	user, err := s.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetProfile(ctx context.Context, id string) (*UserEntity, *shared_errors.AppError) {
	user, err := s.userRepository.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
