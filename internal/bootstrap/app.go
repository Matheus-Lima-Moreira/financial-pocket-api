package bootstrap

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/auth"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/config"
	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/database"
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

	// Auth
	if err := auth.Migrate(db); err != nil {
		return nil, err
	}
	authRepository := auth.NewGormRepository(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	authService := auth.NewService(authRepository, jwtManager)
	authHandler := auth.NewHandler(authService)

	// Router
	router := server.NewRouter(server.Dependencies{
		Logger:      logger,
		AuthHandler: authHandler,
		JWTSecret:   cfg.JWTSecret,
	})

	return &App{
		Router: router,
		Config: cfg,
		Logger: logger,
		DB:     db,
	}, nil
}
