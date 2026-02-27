package action

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/security"
	"gorm.io/gorm/clause"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	models := make([]ActionSchema, 0, len(security.ActionCatalog))
	for _, item := range security.ActionCatalog {
		models = append(models, ActionSchema{
			ID:          item.ID,
			Resource:    item.Resource,
			Action:      item.Action,
			Label:       item.Label,
			Description: item.Description,
		})
	}

	if len(models) == 0 {
		return nil
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"resource", "action", "label", "description"}),
	}).Create(&models).Error
}
