package user

func toModel(user *UserEntity) *UserSchema {
	return &UserSchema{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}
}

func toDomain(model *UserSchema) *UserEntity {
	return &UserEntity{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
	}
}
