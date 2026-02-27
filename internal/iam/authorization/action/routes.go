package action

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/security"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler, requireAction func(string) gin.HandlerFunc) {
	actions := private.Group("/actions")
	actions.GET("/", requireAction(security.ActionActionsList), handler.List)
}
