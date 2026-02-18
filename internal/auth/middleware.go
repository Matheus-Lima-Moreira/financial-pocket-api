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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []string{shared_errors.NewUnauthorized("missing token").Error()},
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			statusCode := http.StatusUnauthorized
			if err == shared_errors.NewUnauthorized("expired token") {
				statusCode = http.StatusUnauthorized
			}
			c.AbortWithStatusJSON(statusCode, gin.H{
				"errors": []string{err.Error()},
			})
			return
		}

		tokenType, ok := claims["type"].(string)
		if !ok || tokenType != "access" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"errors": []string{shared_errors.NewUnauthorized("invalid token").Error()},
			})
			return
		}

		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}