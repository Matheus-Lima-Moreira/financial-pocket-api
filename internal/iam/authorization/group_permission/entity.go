package group_permission

import "time"

type GroupPermissionEntity struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Type      GroupPermissionType `json:"type"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}
