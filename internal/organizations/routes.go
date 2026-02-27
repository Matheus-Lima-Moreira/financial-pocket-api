package organizations

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/security"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler, requireAction func(string) gin.HandlerFunc) {
	organizations := private.Group("/organizations")
	organizations.GET("/", requireAction(security.ActionOrganizationsList), handler.List)
	organizations.GET("/:id", requireAction(security.ActionOrganizationsDetails), handler.Details)
	organizations.POST("/", requireAction(security.ActionOrganizationsCreate), handler.Create)
	organizations.PUT("/:id", requireAction(security.ActionOrganizationsUpdate), handler.Update)
	organizations.DELETE("/:id", requireAction(security.ActionOrganizationsDelete), handler.Delete)
}
