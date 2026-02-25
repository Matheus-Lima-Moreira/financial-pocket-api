package token

import (
	"context"

	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Repository interface {
	Create(ctx context.Context, token *TokenEntity) *shared_errors.AppError
	FindByToken(ctx context.Context, token string) (*TokenEntity, *shared_errors.AppError)
	UpdateStatus(ctx context.Context, token string, status TokenStatus) *shared_errors.AppError
	Delete(ctx context.Context, token string) *shared_errors.AppError
}
