package group_permission

import (
	"time"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/action"
	organization "github.com/Matheus-Lima-Moreira/financial-pocket/internal/organizations"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupPermissionSchema struct {
	ID             string              `gorm:"primaryKey;type:char(36)"`
	Name           string              `gorm:"not null;size:255"`
	Type           GroupPermissionType `gorm:"not null;size:255"`
	OrganizationID string              `gorm:"not null;type:char(36)"`
	CreatedAt      time.Time           `gorm:"autoCreateTime"`
	UpdatedAt      time.Time           `gorm:"autoUpdateTime"`

	GroupPermissionActions []GroupPermissionActionSchema   `gorm:"foreignKey:GroupPermissionID"`
	Organization           organization.OrganizationSchema `gorm:"foreignKey:OrganizationID"`
}

func (GroupPermissionSchema) TableName() string {
	return "group_permissions"
}

func (g *GroupPermissionSchema) BeforeCreate(tx *gorm.DB) error {
	g.ID = uuid.New().String()
	return nil
}

type GroupPermissionActionSchema struct {
	GroupPermissionID string `gorm:"not null;type:char(36);index:idx_group_permission_action,unique"`
	ActionID          string `gorm:"not null;size:255;index:idx_group_permission_action,unique"`

	GroupPermission GroupPermissionSchema `gorm:"foreignKey:GroupPermissionID"`
	Action          action.ActionSchema   `gorm:"foreignKey:ActionID"`
}

func (GroupPermissionActionSchema) TableName() string {
	return "group_permission_actions"
}
