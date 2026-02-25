package user

import (
	"context"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Repository interface {
	Create(ctx context.Context, user *UserEntity) *shared_errors.AppError
	FindByEmail(ctx context.Context, email string) (*UserEntity, *shared_errors.AppError)
	SetEmailVerified(ctx context.Context, id uint, value bool) *shared_errors.AppError
	List(ctx context.Context, page int) ([]UserEntity, *dtos.PaginationDTO, *shared_errors.AppError)
	GetById(ctx context.Context, id uint) (*UserEntity, *shared_errors.AppError)
}
