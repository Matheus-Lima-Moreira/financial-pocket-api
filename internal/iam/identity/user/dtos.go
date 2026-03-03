package user

import "time"

type ListRequest struct {
	Page int `form:"page" binding:"required,min=1"`
}

type DetailsRequest struct {
	ID string `uri:"id" binding:"required,min=1"`
}

type UserReplyDto struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Active         bool         `json:"active"`
	IsPrimary      bool         `json:"is_primary"`
	OrganizationID string       `json:"organization_id"`
	AvatarUrl      string       `json:"avatar_url"`
	RegisterFrom   RegisterFrom `json:"register_from"`
	EmailVerified  bool         `json:"email_verified"`
	State          UserState    `json:"state"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
