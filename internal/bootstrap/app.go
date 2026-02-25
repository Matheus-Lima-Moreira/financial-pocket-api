package bootstrap

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/config"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/database"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/identity/auth"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/identity/user"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/server"
)

type App struct {
	Router *gin.Engine
	Config *config.Config
	Logger zerolog.Logger
	DB     *gorm.DB
}

func NewApp() (*App, error) {
	// Config
	cfg := config.Load()

	// Logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Timestamp().
		Logger()

	// ===== Infrastructure wiring =====
	db, err := database.NewMySQL(cfg.DSN())
	if err != nil {
		return nil, err
	}

	userRepository := user.NewGormRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	jwtManager := auth.NewJWTManager(cfg.AccessTokenSecret, cfg.RefreshTokenSecret)
	authService := auth.NewService(userRepository, jwtManager)
	authHandler := auth.NewHandler(authService)

	// Router
	router := server.NewRouter(server.Dependencies{
		Logger:      logger,
		Config:      cfg,
		JWTManager:  jwtManager,
		AuthHandler: authHandler,
		UserHandler: userHandler,
	})

	return &App{
		Router: router,
		Config: cfg,
		Logger: logger,
		DB:     db,
	}, nil
}
