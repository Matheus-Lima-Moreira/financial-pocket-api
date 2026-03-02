package group_permission

import (
	"context"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/consts"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context, page int) ([]GroupPermissionEntity, *dtos.PaginationDTO, *shared_errors.AppError) {
	var models []GroupPermissionSchema

	limit := consts.PaginationDefaultLimit

	err := r.db.WithContext(ctx).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&models).Error

	if err != nil {
		return nil, nil, shared_errors.NewBadRequest(err.Error())
	}

	domains := make([]GroupPermissionEntity, len(models))
	for i, model := range models {
		domains[i] = *toDomain(&model)
	}

	var total int64 = 0

	err = r.db.WithContext(ctx).
		Model(&GroupPermissionSchema{}).
		Count(&total).Error

	if err != nil {
		return nil, nil, shared_errors.NewBadRequest(err.Error())
	}

	next := 0
	previous := 0
	if page*limit < int(total) {
		next = page + 1
	}
	if page > 1 {
		previous = page - 1
	}

	pagination := &dtos.PaginationDTO{
		Page:     page,
		Limit:    limit,
		Total:    int(total),
		Next:     next,
		Previous: previous,
	}

	return domains, pagination, nil
}

func (r *GormRepository) Details(ctx context.Context, id uint) (*GroupPermissionEntity, *shared_errors.AppError) {
	var model GroupPermissionSchema

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return toDomain(&model), nil
}

func (r *GormRepository) Create(ctx context.Context, groupPermission *GroupPermissionEntity) *shared_errors.AppError {
	model := toModel(groupPermission)

	err := r.db.WithContext(ctx).
		Create(model).Error

	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}

func (r *GormRepository) Update(ctx context.Context, groupPermission *GroupPermissionEntity) *shared_errors.AppError {
	model := toModel(groupPermission)

	err := r.db.WithContext(ctx).
		Model(&GroupPermissionSchema{}).
		Where("id = ?", groupPermission.ID).
		Updates(model).Error

	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id uint) *shared_errors.AppError {
	err := r.db.WithContext(ctx).
		Model(&GroupPermissionSchema{}).
		Where("id = ?", id).
		Delete(&GroupPermissionSchema{}).Error

	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}

func (r *GormRepository) GetAllOfTypeSystem(ctx context.Context) ([]GroupPermissionEntity, *shared_errors.AppError) {
	var models []GroupPermissionSchema

	err := r.db.WithContext(ctx).
		Where("type = ?", GroupPermissionSystem).
		Find(&models).Error

	if err != nil {
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	domains := make([]GroupPermissionEntity, len(models))
	for i, model := range models {
		domains[i] = *toDomain(&model)
	}

	return domains, nil
}
