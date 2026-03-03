package user

import (
	"context"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/dtos"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

type Service struct {
	userRepository Repository
}

func NewService(userRepository Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) List(ctx context.Context, page int, organizationID string) ([]UserReplyDto, *dtos.PaginationDTO, *shared_errors.AppError) {
	users, pagination, err := s.userRepository.List(ctx, page, organizationID)
	if err != nil {
		return nil, nil, shared_errors.NewBadRequest(err.Error())
	}

	usersDTO := make([]UserReplyDto, len(users))
	for i, user := range users {
		usersDTO[i] = UserReplyDto{
			ID:             user.ID,
			Name:           user.Name,
			Email:          user.Email,
			IsPrimary:      user.IsPrimary,
			Active:         user.Active,
			OrganizationID: user.OrganizationID,
			AvatarUrl:      user.Avatar,
			RegisterFrom:   user.RegisterFrom,
			EmailVerified:  user.EmailVerified,
			State:          DetermineUserState(&user),
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		}
	}
	return usersDTO, pagination, nil
}

func DetermineUserState(user *UserEntity) UserState {
	if !user.Active {
		return UserStateInactive
	}
	if user.RegisterFrom == RegisterFromInvite && !user.EmailVerified {
		return UserStateInvitePending
	}
	return UserStateActive
}

func (s *Service) Details(ctx context.Context, id string) (*UserEntity, *shared_errors.AppError) {
	user, err := s.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetProfile(ctx context.Context, id string) (*UserEntity, *shared_errors.AppError) {
	user, err := s.userRepository.GetProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
