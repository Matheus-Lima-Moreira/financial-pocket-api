package auth

import (
	"context"
	"net/url"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/identity/user"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Service struct {
	userRepository user.Repository
	jwt            *JWTManager
	emailSender    EmailSender
	verifyBaseURL  string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewService(userRepository user.Repository, jwt *JWTManager, emailSender EmailSender, verifyBaseURL string) *Service {
	return &Service{
		userRepository: userRepository,
		jwt:            jwt,
		emailSender:    emailSender,
		verifyBaseURL:  verifyBaseURL,
	}
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Service) Register(ctx context.Context, input RegisterInput) *shared_errors.AppError {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	user := &user.UserEntity{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hash),
	}

	existingUser, repositoryErr := s.userRepository.FindByEmail(ctx, input.Email)
	if repositoryErr != nil && repositoryErr.ErrorCode != shared_errors.CodeNotFound {
		return repositoryErr
	}

	if existingUser != nil {
		return shared_errors.NewConflict("email already in use", "email")
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return err
	}

	if err := s.sendVerificationEmail(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPair, *shared_errors.AppError) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !user.EmailVerified {
		return nil, shared_errors.NewUnauthorized("auth.email_not_verified")
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

func (s *Service) VerifyEmail(ctx context.Context, verifyToken string) *shared_errors.AppError {
	userID, err := s.jwt.ValidateVerifyEmailToken(verifyToken)
	if err != nil {
		return err
	}

	parsedID, parseErr := strconv.Atoi(userID)
	if parseErr != nil {
		return shared_errors.NewUnauthorized("invalid token")
	}

	if repositoryErr := s.userRepository.SetEmailVerified(ctx, uint(parsedID), true); repositoryErr != nil {
		return repositoryErr
	}

	return nil
}

func (s *Service) sendVerificationEmail(ctx context.Context, createdUser *user.UserEntity) *shared_errors.AppError {
	userID := strconv.Itoa(int(createdUser.ID))
	verifyToken, tokenErr := s.jwt.GenerateVerifyEmailToken(userID)
	if tokenErr != nil {
		return tokenErr
	}

	verifyURL := s.verifyBaseURL + "?token=" + url.QueryEscape(verifyToken)
	if err := s.emailSender.SendVerifyEmail(ctx, createdUser.Email, createdUser.Name, verifyURL); err != nil {
		return shared_errors.NewBadRequest("auth.verify_email_send_failed")
	}

	return nil
}
