package token

import (
	"context"
	"errors"

	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func isTokenAlreadyExistsError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrDuplicatedKey)
}

func (r *GormRepository) Create(ctx context.Context, token *TokenEntity) *shared_errors.AppError {
	model := toModel(token)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isTokenAlreadyExistsError(err) {
			return shared_errors.NewConflict("token already exists", "token")
		}
		return shared_errors.NewBadRequest(err.Error())
	}

	*token = *toDomain(model)

	return nil
}

func (r *GormRepository) FindByToken(ctx context.Context, token string) (*TokenEntity, *shared_errors.AppError) {
	var model TokenSchema

	err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.NewNotFound("token")
		}
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return toDomain(&model), nil
}

func (r *GormRepository) UpdateStatus(ctx context.Context, token string, status TokenStatus) *shared_errors.AppError {
	if err := r.db.WithContext(ctx).
		Model(&TokenSchema{}).
		Where("token = ?", token).
		Update("status", status).Error; err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}

func (r *GormRepository) Delete(ctx context.Context, token string) *shared_errors.AppError {
	if err := r.db.WithContext(ctx).
		Model(&TokenSchema{}).
		Where("token = ?", token).
		Delete(&TokenSchema{}).Error; err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}
	return nil
}
