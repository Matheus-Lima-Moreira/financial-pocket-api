package user

func toModel(user *UserEntity) *UserSchema {
	return &UserSchema{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Password:      user.Password,
		EmailVerified: user.EmailVerified,
		Avatar:        user.Avatar,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

func toDomain(model *UserSchema) *UserEntity {
	return &UserEntity{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		Password:      model.Password,
		EmailVerified: model.EmailVerified,
		Avatar:        model.Avatar,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
