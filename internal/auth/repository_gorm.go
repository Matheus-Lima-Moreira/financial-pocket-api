package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

type userModel struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null;size:255"`
	Password  string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (r *GormRepository) Create(ctx context.Context, user *User) error {
	model := toModel(user)

	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUserAlreadyExists
		}
		return err
	}

	*user = *toDomain(model)

	return nil
}

func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var model userModel

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return toDomain(&model), nil
}

func toModel(user *User) *userModel {
	return &userModel{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

func toDomain(model *userModel) *User {
	return &User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}
}
