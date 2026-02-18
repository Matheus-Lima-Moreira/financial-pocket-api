package auth

import (
	"context"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Service struct {
	repo Repository
	jwt  *JWTManager
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewService(repo Repository, jwt *JWTManager) *Service {
	return &Service{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) *shared_errors.AppError {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	user := &User{
		Email:    email,
		Password: string(hash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, *shared_errors.AppError) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, shared_errors.NewUnauthorized("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, shared_errors.NewUnauthorized("invalid credentials")
	}

	userID := strconv.Itoa(int(user.ID))

	accessToken, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, *shared_errors.AppError) {
	userID, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, shared_errors.NewUnauthorized("invalid token")
	}

	accessToken, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, shared_errors.NewBadRequest(err.Error())
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
