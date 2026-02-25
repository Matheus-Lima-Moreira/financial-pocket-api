package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
)

func AuthMiddleware(jwtManager *JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			errorDetail := shared_errors.NewUnauthorized("missing token").ToErrorDetail()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []shared_errors.ErrorDetail{errorDetail},
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, appErr := jwtManager.ValidateToken(tokenStr)
		if appErr != nil {
			errorDetail := appErr.ToErrorDetail()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []shared_errors.ErrorDetail{errorDetail},
			})
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			errorDetail := shared_errors.NewUnauthorized("invalid token").ToErrorDetail()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []shared_errors.ErrorDetail{errorDetail},
			})
			return
		}

		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}