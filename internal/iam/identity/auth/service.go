package auth

import (
	"context"
	"net/url"

	"golang.org/x/crypto/bcrypt"

	group_permission "github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/group_permission"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/identity/user"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/provisioning/token"
	emails "github.com/Matheus-Lima-Moreira/financial-pocket/internal/notifications/emails"
	organization "github.com/Matheus-Lima-Moreira/financial-pocket/internal/organizations"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Service struct {
	userRepository            user.Repository
	organizationRepository    organization.Repository
	groupPermissionRepository group_permission.Repository
	jwt                       *JWTManager
	tokenService              *token.Service
	emailSender               emails.EmailSender
	frontendBaseURL           string
}

func NewService(userRepository user.Repository, organizationRepository organization.Repository, groupPermissionRepository group_permission.Repository, jwt *JWTManager, tokenService *token.Service, emailSender emails.EmailSender, frontendBaseURL string) *Service {
	return &Service{
		userRepository:            userRepository,
		organizationRepository:    organizationRepository,
		groupPermissionRepository: groupPermissionRepository,
		jwt:                       jwt,
		tokenService:              tokenService,
		emailSender:               emailSender,
		frontendBaseURL:           frontendBaseURL,
	}
}

func (s *Service) Register(ctx context.Context, input RegisterInputDTO) *shared_errors.AppError {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return shared_errors.NewBadRequest(err.Error())
	}

	user := &user.UserEntity{
		Name:         input.User.Name,
		Email:        input.User.Email,
		Password:     string(hash),
		RegisterFrom: user.RegisterFromForm,
		IsPrimary:    true,
	}

	existingUser, repositoryErr := s.userRepository.FindByEmail(ctx, input.User.Email)
	if repositoryErr != nil && repositoryErr.ErrorCode != shared_errors.CodeNotFound {
		return repositoryErr
	}

	if existingUser != nil {
		return nil
	}

	organization := &organization.OrganizationEntity{
		Cellphone: input.Organization.Cellphone,
		Name:      input.Organization.Name,
	}

	if err := s.organizationRepository.Create(ctx, organization); err != nil {
		return err
	}

	user.OrganizationID = organization.ID

	if err := s.userRepository.Create(ctx, user); err != nil {
		return err
	}

	if user.IsPrimary {
		groupPermissions, err := s.groupPermissionRepository.GetAllOfTypeSystem(ctx)
		if err != nil {
			return err
		}

		for _, groupPermission := range groupPermissions {
			if err := s.userRepository.AddGroupPermission(ctx, user.ID, groupPermission.ID); err != nil {
				return err
			}
		}
	}

	if err := s.SendVerificationEmail(ctx, user.Email); err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*TokenPairDTO, *shared_errors.AppError) {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, shared_errors.NewUnauthorized("error.invalid_credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return nil, shared_errors.NewUnauthorized("error.invalid_credentials")
	}

	accessToken, err := s.jwt.GenerateAccessToken(user.ID, user.OrganizationID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(user.ID, user.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &TokenPairDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPairDTO, *shared_errors.AppError) {
	claims, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID := claims.UserID
	organizationID := claims.OrganizationID

	accessToken, err := s.jwt.GenerateAccessToken(userID, organizationID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(userID, organizationID)
	if err != nil {
		return nil, err
	}

	return &TokenPairDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) ResetPassword(ctx context.Context, resetPasswordToken string, newPassword string) *shared_errors.AppError {
	resource := token.TokenResourceResetPassword
	valid, err := s.tokenService.VerifyToken(ctx, resetPasswordToken, &resource)
	if err != nil {
		return err
	}

	if !valid {
		return shared_errors.NewUnauthorized("error.invalid_token")
	}

	tokenEntity, err := s.tokenService.GetByToken(ctx, resetPasswordToken)
	if err != nil {
		return err
	}

	hash, errHash := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return shared_errors.NewBadRequest(errHash.Error())
	}

	if err := s.userRepository.UpdatePassword(ctx, tokenEntity.ReferenceID, string(hash)); err != nil {
		return err
	}

	if err := s.tokenService.UpdateStatus(ctx, resetPasswordToken, token.TokenStatusUsed); err != nil {
		return err
	}

	return nil
}

func (s *Service) SendVerificationEmail(ctx context.Context, email string) *shared_errors.AppError {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	tokenEntity, err := s.tokenService.GenerateToken(ctx, token.TokenResourceVerifyEmail, user.ID, map[string]any{
		"email": email,
		"name":  user.Name,
	})
	if err != nil {
		return err
	}

	verifyURL := s.frontendBaseURL + "/auth/verify-email?token=" + url.QueryEscape(tokenEntity.Token)

	if err := s.emailSender.SendVerifyEmail(ctx, email, user.Name, verifyURL); err != nil {
		return shared_errors.NewBadRequest("auth.verify_email_send_failed")
	}

	return nil
}

func (s *Service) SendResetPasswordEmail(ctx context.Context, email string) *shared_errors.AppError {
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	tokenEntity, err := s.tokenService.GenerateToken(ctx, token.TokenResourceResetPassword, user.ID, map[string]any{
		"email": email,
		"name":  user.Name,
	})
	if err != nil {
		return err
	}

	resetPasswordURL := s.frontendBaseURL + "/auth/reset-password?token=" + url.QueryEscape(tokenEntity.Token)
	if err := s.emailSender.SendResetPasswordEmail(ctx, email, user.Name, resetPasswordURL); err != nil {
		return shared_errors.NewBadRequest("auth.reset_password_send_failed")
	}

	return nil
}

func (s *Service) VerifyEmail(ctx context.Context, verifyToken string) *shared_errors.AppError {
	resource := token.TokenResourceVerifyEmail
	valid, err := s.tokenService.VerifyToken(ctx, verifyToken, &resource)
	if err != nil {
		return err
	}

	if !valid {
		return shared_errors.NewUnauthorized("error.invalid_token")
	}

	tokenEntity, err := s.tokenService.GetByToken(ctx, verifyToken)
	if err != nil {
		return err
	}

	if err := s.tokenService.UpdateStatus(ctx, verifyToken, token.TokenStatusUsed); err != nil {
		return err
	}

	if err := s.userRepository.SetEmailVerified(ctx, tokenEntity.ReferenceID, true); err != nil {
		return err
	}

	return nil
}
