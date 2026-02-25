package token

import "time"

type TokenResource string

const (
	TokenResourceVerifyEmail   TokenResource = "users:verify-email"
	TokenResourceResetPassword TokenResource = "users:reset-password"
)

type TokenStatus string

const (
	TokenStatusActive   TokenStatus = "ACTIVE"
	TokenStatusInactive TokenStatus = "INACTIVE"
	TokenStatusExpired  TokenStatus = "EXPIRED"
	TokenStatusUsed     TokenStatus = "USED"
)

type TokenEntity struct {
	ID          uint           `json:"id"`
	ReferenceID string         `json:"reference_id"`
	Token       string         `json:"token"`
	Resource    TokenResource  `json:"resource"`
	ExpiresAt   time.Time      `json:"expires_at"`
	Status      TokenStatus    `json:"status"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func NewTokenEntity(resource TokenResource, referenceID string, metadata map[string]any) *TokenEntity {
	return &TokenEntity{
		ReferenceID: referenceID,
		Resource:    resource,
		Token:       generateToken(),
		ExpiresAt:   time.Now().Add(time.Hour * 24),
		Status:      TokenStatusActive,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
