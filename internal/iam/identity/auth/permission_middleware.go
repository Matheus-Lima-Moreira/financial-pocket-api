package auth

import (
	"context"
	"net/http"

	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionMiddleware struct {
	db *gorm.DB
}

func NewPermissionMiddleware(db *gorm.DB) *PermissionMiddleware {
	return &PermissionMiddleware{db: db}
}

func (m *PermissionMiddleware) Require(actionID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDAny, exists := c.Get("user_id")
		if !exists {
			errorDetail := shared_errors.NewUnauthorized("error.invalid_token").ToErrorDetail()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []shared_errors.ErrorDetail{errorDetail},
			})
			return
		}

		userID, ok := userIDAny.(string)
		if !ok || userID == "" {
			errorDetail := shared_errors.NewUnauthorized("error.invalid_token").ToErrorDetail()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []shared_errors.ErrorDetail{errorDetail},
			})
			return
		}

		hasPermission, err := m.hasPermission(c.Request.Context(), userID, actionID)
		if err != nil {
			c.Error(shared_errors.NewBadRequest(err.Error()))
			c.Abort()
			return
		}

		if !hasPermission {
			errorDetail := shared_errors.NewUnauthorized("auth.not_allowed").ToErrorDetail()
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"errors": []shared_errors.ErrorDetail{errorDetail},
			})
			return
		}

		c.Next()
	}
}

func (m *PermissionMiddleware) hasPermission(ctx context.Context, userID, actionID string) (bool, error) {
	var total int64
	err := m.db.WithContext(ctx).
		Table("user_group_permissions AS ugp").
		Joins("JOIN group_permission_actions AS gpa ON gpa.group_permission_id = ugp.group_permission_id").
		Where("ugp.user_id = ? AND gpa.action_id = ?", userID, actionID).
		Count(&total).Error
	if err != nil {
		return false, err
	}

	return total > 0, nil
}
