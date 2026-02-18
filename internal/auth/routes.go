package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {
	group := router.Group("/auth")

	group.POST("/register", handler.Register)
	group.POST("/login", handler.Login)
}
