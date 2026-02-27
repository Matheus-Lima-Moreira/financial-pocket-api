package group_permission

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/action"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/organizations"
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

	groupPermission := GroupPermissionSchema{
		Name:           "System",
		Type:           GroupPermissionSystem,
		OrganizationID: systemOrganization.ID,
	}

	if err := db.
		Where("name = ? AND type = ? AND organization_id = ?", groupPermission.Name, groupPermission.Type, groupPermission.OrganizationID).
		FirstOrCreate(&groupPermission).Error; err != nil {
		return err
	}

	var actions []action.ActionSchema
	if err := db.Find(&actions).Error; err != nil {
		return err
	}

	if len(actions) == 0 {
		return nil
	}

	associations := make([]GroupPermissionActionSchema, 0, len(actions))
	for _, item := range actions {
		associations = append(associations, GroupPermissionActionSchema{
			GroupPermissionID: groupPermission.ID,
			ActionID:          item.ID,
		})
	}

	return db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "group_permission_id"}, {Name: "action_id"}},
			DoNothing: true,
		}).
		Create(&associations).Error
}
