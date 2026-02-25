package token

import "time"

type TokenSchema struct {
	ID          uint           `gorm:"primaryKey"`
	Token       string         `gorm:"uniqueIndex;not null;size:255"`
	ReferenceID string         `gorm:"not null;size:255"`
	Resource    TokenResource  `gorm:"not null"`
	ExpiresAt   time.Time      `gorm:"not null"`
	Status      TokenStatus    `gorm:"not null;default:ACTIVE"`
	Metadata    map[string]any `gorm:"type:json"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

func (TokenSchema) TableName() string {
	return "tokens"
}
