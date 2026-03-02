package user

import (
	"context"
	"errors"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/consts"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const mysqlErrDuplicateEntry = 1062

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == mysqlErrDuplicateEntry
	}
	return false
}

func (r *GormRepository) Create(ctx context.Context, user *UserEntity) *shared_errors.AppError {
	model := toModel(user)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return shared_errors.NewConflict("error.email_already_in_use", "email")
		}
		return shared_errors.NewBadRequest(err.Error())
	}

	*user = *toDomain(model)

	return nil
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*UserEntity, *shared_errors.AppError) {
	var model UserSchema

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.NewNotFound("user")
		}
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return toDomain(&model), nil
}

func (r *GormRepository) SetEmailVerified(ctx context.Context, id string, value bool) *shared_errors.AppError {
	result := r.db.WithContext(ctx).
		Model(&UserSchema{}).
		Where("id = ?", id).
		Update("email_verified", value)

	if result.Error != nil {
		return shared_errors.NewBadRequest(result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return shared_errors.NewNotFound("user")
	}

	return nil
}

func (r *GormRepository) List(ctx context.Context, page int) ([]UserEntity, *dtos.PaginationDTO, *shared_errors.AppError) {
	var models []UserSchema

	limit := consts.PaginationDefaultLimit

	err := r.db.WithContext(ctx).
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&models).Error

	if err != nil {
		return nil, nil, shared_errors.NewBadRequest(err.Error())
	}

	domains := make([]UserEntity, len(models))
	for i, model := range models {
		domains[i] = *toDomain(&model)
	}

	var total int64 = 0

	err = r.db.WithContext(ctx).
		Model(&UserSchema{}).
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

func (r *GormRepository) GetById(ctx context.Context, id string) (*UserEntity, *shared_errors.AppError) {
	var model UserSchema

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.NewNotFound("user")
		}
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return toDomain(&model), nil
}

func (r *GormRepository) UpdatePassword(ctx context.Context, id string, password string) *shared_errors.AppError {
	err := r.db.WithContext(ctx).
		Model(&UserSchema{}).
		Where("id = ?", id).
		Update("password", password).Error

	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}

func (r *GormRepository) GetProfile(ctx context.Context, id string) (*UserEntity, *shared_errors.AppError) {
	var model UserSchema

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.NewNotFound("user")
		}
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return toDomain(&model), nil
}

func (r *GormRepository) AddGroupPermission(ctx context.Context, userID string, groupPermissionID string) *shared_errors.AppError {
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("group_permission_id = ?", groupPermissionID).
		FirstOrCreate(&UserGroupPermissionSchema{
			UserID:            userID,
			GroupPermissionID: groupPermissionID,
		}).Error

	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}
