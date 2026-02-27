package group_permission

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/security"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler, requireAction func(string) gin.HandlerFunc) {
	groupPermissions := private.Group("/group-permissions")
	groupPermissions.GET("/", requireAction(security.ActionGroupPermissionsList), handler.List)
	groupPermissions.GET("/:id", requireAction(security.ActionGroupPermissionsDetails), handler.Details)
	groupPermissions.POST("/", requireAction(security.ActionGroupPermissionsCreate), handler.Create)
	groupPermissions.PUT("/:id", requireAction(security.ActionGroupPermissionsUpdate), handler.Update)
	groupPermissions.DELETE("/:id", requireAction(security.ActionGroupPermissionsDelete), handler.Delete)
}
