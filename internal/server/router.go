package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/config"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/identity/auth"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/identity/user"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/middlewares"
)

type Dependencies struct {
	Logger      zerolog.Logger
	Config      *config.Config
	JWTManager  *auth.JWTManager
	AuthHandler *auth.Handler
	UserHandler *user.Handler
}

func NewRouter(dep Dependencies) *gin.Engine {
	if dep.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	if dep.Config.TrustedProxies() != nil {
		router.SetTrustedProxies(dep.Config.TrustedProxies())
	} else {
		router.SetTrustedProxies(nil)
	}

	router.Use(gin.Recovery())
	router.Use(middlewares.LoggerMiddleware(dep.Logger))
	router.Use(middlewares.ErrorMiddleware())

	public := router.Group("")
	private := router.Group("")

	// Middlewares
	private.Use(auth.AuthMiddleware(dep.JWTManager))

	// Routes
	auth.RegisterRoutes(public, private, dep.AuthHandler)
	user.RegisterRoutes(public, private, dep.UserHandler)

	public.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})

	return router
}
