package token

import "github.com/gin-gonic/gin"

func RegisterRoutes(public, private *gin.RouterGroup, handler *Handler) {
	privateTokensRoute := private.Group("/tokens")
	privateTokensRoute.POST("/resend-verification-email", handler.ResendVerificationEmail)

	publicTokensRoute := public.Group("/tokens")
	publicTokensRoute.POST("/reset-password", handler.ResetPassword)
	publicTokensRoute.GET("/verify-email", handler.VerifyEmail)
}
