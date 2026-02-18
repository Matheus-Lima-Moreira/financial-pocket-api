package auth

import (
	"context"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Repository interface {
	Create(ctx context.Context, user *User) *shared_errors.AppError
	FindByEmail(ctx context.Context, email string) (*User, *shared_errors.AppError)
}
