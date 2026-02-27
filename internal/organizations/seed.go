package organizations

import "gorm.io/gorm"

const (
	SystemOrganizationName      = "System Organization"
	SystemOrganizationCellphone = "+5500000000000"
	SystemOrganizationLogo      = "https://financial-pocket.dev/assets/system-organization-logo.png"
)

func Seed(db *gorm.DB) error {
	return db.
		Where("name = ?", SystemOrganizationName).
		FirstOrCreate(&OrganizationSchema{
			Name:      SystemOrganizationName,
			Cellphone: SystemOrganizationCellphone,
			Logo:      SystemOrganizationLogo,
		}).Error
}
