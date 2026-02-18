package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/auth"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/middlewares"
)

type Dependencies struct {
	Logger      zerolog.Logger
	AuthHandler *auth.Handler
	JWTSecret   string
}

func NewRouter(dep Dependencies) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(dep.Logger))
	router.Use(middlewares.ErrorMiddleware())
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	auth.RegisterRoutes(router, dep.AuthHandler)

	protected := router.Group("/api")
	protected.Use(auth.AuthMiddleware(dep.JWTSecret))

	return router
}
