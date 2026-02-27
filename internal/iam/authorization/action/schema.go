package action

import (
	"fmt"
	"strings"
	"gorm.io/gorm"
)

type ActionSchema struct {
	ID          string `gorm:"primaryKey;size:255"`
	Resource    string `gorm:"not null;size:255"`
	Action      string `gorm:"not null;size:255"`
	Label       string `gorm:"not null;size:255"`
	Description string `gorm:"not null;size:255"`
}

func (ActionSchema) TableName() string {
	return "actions"
}

func (a *ActionSchema) BeforeCreate(tx *gorm.DB) error {
	a.Resource = strings.ToLower(strings.TrimSpace(a.Resource))
	a.Action = strings.ToLower(strings.TrimSpace(a.Action))
	a.ID = fmt.Sprintf("%s:%s", a.Resource, a.Action)
	return nil
}
