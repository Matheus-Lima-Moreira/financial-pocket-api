package group_permission

import (
	"context"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Repository interface {
	List(ctx context.Context, page int) ([]GroupPermissionEntity, *dtos.PaginationDTO, *shared_errors.AppError)
	Details(ctx context.Context, id uint) (*GroupPermissionEntity, *shared_errors.AppError)
	Create(ctx context.Context, groupPermission *GroupPermissionEntity) *shared_errors.AppError
	Update(ctx context.Context, groupPermission *GroupPermissionEntity) *shared_errors.AppError
	Delete(ctx context.Context, id uint) *shared_errors.AppError
	GetAllOfTypeSystem(ctx context.Context) ([]GroupPermissionEntity, *shared_errors.AppError)
}
