package user

import "time"

type UserSchema struct {
	ID            uint      `gorm:"primaryKey"`
	Name          string    `gorm:"not null;size:255"`
	Email         string    `gorm:"uniqueIndex;not null;size:255"`
	Password      string    `gorm:"not null;size:255"`
	EmailVerified bool      `gorm:"not null;default:false"`
	Avatar        string    `gorm:"not null;size:255"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (UserSchema) TableName() string {
	return "users"
}
