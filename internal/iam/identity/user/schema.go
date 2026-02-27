package user

import (
	"time"

	group_permission "github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/group_permission"
	organization "github.com/Matheus-Lima-Moreira/financial-pocket/internal/organizations"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSchema struct {
	ID             string    `gorm:"primaryKey;type:char(36)"`
	Name           string    `gorm:"not null;size:255"`
	Email          string    `gorm:"uniqueIndex;not null;size:255"`
	Password       string    `gorm:"not null;size:255"`
	Active         bool      `gorm:"not null;default:true"`
	IsPrimary      bool      `gorm:"not null;default:false"`
	OrganizationID string    `gorm:"not null;type:char(36)"`
	Avatar         string    `gorm:"not null;size:255"`
	RegisterFrom   string    `gorm:"not null;size:255"`
	EmailVerified  bool      `gorm:"not null;default:false"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`

	UserGroupPermissions []UserGroupPermissionSchema     `gorm:"foreignKey:UserID"`
	Organization         organization.OrganizationSchema `gorm:"foreignKey:OrganizationID"`
}

func (UserSchema) TableName() string {
	return "users"
}

func (u *UserSchema) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	return nil
}

type UserGroupPermissionSchema struct {
	UserID            string `gorm:"not null;type:char(36);index:idx_user_group_permission,unique"`
	GroupPermissionID string `gorm:"not null;type:char(36);index:idx_user_group_permission,unique"`

	User            UserSchema                             `gorm:"foreignKey:UserID"`
	GroupPermission group_permission.GroupPermissionSchema `gorm:"foreignKey:GroupPermissionID"`
}

func (UserGroupPermissionSchema) TableName() string {
	return "user_group_permissions"
}
