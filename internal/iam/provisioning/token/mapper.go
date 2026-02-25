package token

func toModel(token *TokenEntity) *TokenSchema {
	return &TokenSchema{
		ID:        token.ID,
		Token:     token.Token,
		Resource:  token.Resource,
		ExpiresAt: token.ExpiresAt,
		Status:    token.Status,
		CreatedAt: token.CreatedAt,
		UpdatedAt: token.UpdatedAt,
	}
}

func toDomain(model *TokenSchema) *TokenEntity {
	return &TokenEntity{
		ID:        model.ID,
		Token:     model.Token,
		Resource:  model.Resource,
		ExpiresAt: model.ExpiresAt,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
