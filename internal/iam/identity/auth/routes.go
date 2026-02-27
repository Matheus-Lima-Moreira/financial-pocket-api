package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler) {
	publicAuthRoute := public.Group("/auth")

	publicAuthRoute.POST("/register", handler.Register)
	publicAuthRoute.POST("/login", handler.Login)
	publicAuthRoute.POST("/refresh", handler.Refresh)
	publicAuthRoute.POST("/send-reset-password-email", handler.SendResetPasswordEmail)
	publicAuthRoute.POST("/reset-password", handler.ResetPassword)
	publicAuthRoute.GET("/verify-email", handler.VerifyEmail)

	privateAuthRoute := private.Group("/auth")
	privateAuthRoute.POST("/resend-verification-email", handler.ResendVerificationEmail)
}
