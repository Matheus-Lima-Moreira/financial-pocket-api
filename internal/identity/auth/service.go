package auth

import (
	"context"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/identity/user"
)

type Service struct {
	userRepository user.Repository
	jwt  *JWTManager
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewService(userRepository user.Repository, jwt *JWTManager) *Service {
	return &Service{
		userRepository: userRepository,
		jwt:  jwt,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) *shared_errors.AppError {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	user := &user.UserEntity{
		Email:    email,
		Password: string(hash),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, *shared_errors.AppError) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, *shared_errors.AppError) {
	userID, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
