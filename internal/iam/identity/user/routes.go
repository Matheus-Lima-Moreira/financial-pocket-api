package user

import (
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/security"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler, requireAction func(string) gin.HandlerFunc) {
	users := private.Group("/users")
	users.GET("/", requireAction(security.ActionUsersList), handler.List)
	users.GET("/:id", requireAction(security.ActionUsersDetails), handler.Details)
}
