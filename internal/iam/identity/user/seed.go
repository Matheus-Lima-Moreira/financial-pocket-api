package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	password := "@Admin2026"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.
		Where("email = ?", "admin@financial-pocket.dev").
		FirstOrCreate(&UserSchema{
			Name:          "Admin",
			Email:         "admin@financial-pocket.dev",
			Password:      string(hash),
			EmailVerified: true,
		}).Error
}
