package organizations

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationSchema struct {
	ID        string    `gorm:"primaryKey;type:char(36)"`
	Name      string    `gorm:"not null;size:255"`
	Cellphone string    `gorm:"not null;size:255"`
	Logo      string    `gorm:"not null;size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (OrganizationSchema) TableName() string {
	return "organizations"
}

func (o *OrganizationSchema) BeforeCreate(tx *gorm.DB) error {
	o.ID = uuid.New().String()
	return nil
}
