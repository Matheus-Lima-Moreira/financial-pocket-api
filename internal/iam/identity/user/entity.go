package user

import "time"

type UserEntity struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Password       string       `json:"password"`
	EmailVerified  bool         `json:"email_verified"`
	Avatar         string       `json:"avatar"`
	OrganizationID string       `json:"organization_id"`
	RegisterFrom   RegisterFrom `json:"register_from"`
	IsPrimary      bool         `json:"is_primary"`
	Active         bool         `json:"active"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
