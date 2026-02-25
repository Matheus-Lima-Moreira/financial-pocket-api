package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler) {
	users := private.Group("/users")
	users.GET("", handler.List)
	users.GET("/:id", handler.Details)
}
