package auth

import (
	"context"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
	jwt  *JWTManager
}

func NewService(repo Repository, jwt *JWTManager) *Service {
	return &Service{
		repo: repo,
		jwt:  jwt,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Email:    email,
		Password: string(hash),
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.jwt.Generate(strconv.Itoa(int(user.ID)))
}
