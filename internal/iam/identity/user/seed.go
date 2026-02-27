package user

import (
	group_permission "github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/group_permission"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/organizations"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Seed(db *gorm.DB) error {
	var systemOrganization organizations.OrganizationSchema
	if err := db.
		Where("name = ?", organizations.SystemOrganizationName).
		First(&systemOrganization).Error; err != nil {
		return err
	}

	var systemGroupPermission group_permission.GroupPermissionSchema
	if err := db.
		Where("type = ? AND organization_id = ?", group_permission.GroupPermissionSystem, systemOrganization.ID).
		First(&systemGroupPermission).Error; err != nil {
		return err
	}

	password := "@Admin2026"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := UserSchema{
		Name:           "Admin",
		Email:          "admin@financial-pocket.dev",
		Password:       string(hash),
		OrganizationID: systemOrganization.ID,
		Avatar:         "https://financial-pocket.dev/assets/system-user-avatar.png",
		RegisterFrom:   "system-seed",
		EmailVerified:  true,
		Active:         true,
		IsPrimary:      true,
	}

	if err := db.
		Where("email = ?", "admin@financial-pocket.dev").
		FirstOrCreate(&admin).Error; err != nil {
		return err
	}

	association := UserGroupPermissionSchema{
		UserID:            admin.ID,
		GroupPermissionID: systemGroupPermission.ID,
	}

	return db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "group_permission_id"}},
			DoNothing: true,
		}).
		Create(&association).Error
}
