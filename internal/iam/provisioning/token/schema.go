package token

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TokenSchema struct {
	ID          string            `gorm:"primaryKey;type:char(36)"`
	Token       string            `gorm:"uniqueIndex;not null;size:255"`
	ReferenceID string            `gorm:"not null;size:255"`
	Resource    TokenResource     `gorm:"not null"`
	ExpiresAt   time.Time         `gorm:"not null"`
	Status      TokenStatus       `gorm:"not null;default:ACTIVE"`
	Metadata    datatypes.JSONMap `gorm:"type:json"`
	CreatedAt   time.Time         `gorm:"autoCreateTime"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime"`
}

func (TokenSchema) TableName() string {
	return "tokens"
}

func (t *TokenSchema) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New().String()
	return nil
}
