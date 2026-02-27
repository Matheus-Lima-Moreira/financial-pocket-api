package database

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/action"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/authorization/group_permission"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/identity/user"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/iam/provisioning/token"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/organizations"
)

func NewMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Pool configuration
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	if err := RunSeeds(db); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	if err := organizations.Migrate(db); err != nil {
		return err
	}

	if err := action.Migrate(db); err != nil {
		return err
	}

	if err := group_permission.Migrate(db); err != nil {
		return err
	}

	if err := user.Migrate(db); err != nil {
		return err
	}

	if err := token.Migrate(db); err != nil {
		return err
	}

	return nil
}

func RunSeeds(db *gorm.DB) error {
	if err := organizations.Seed(db); err != nil {
		return err
	}

	if err := action.Seed(db); err != nil {
		return err
	}

	if err := group_permission.Seed(db); err != nil {
		return err
	}

	if err := user.Seed(db); err != nil {
		return err
	}

	return nil
}
