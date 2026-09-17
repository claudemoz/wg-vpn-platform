package app

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/device"
	"backend/internal/user"
	vpnserver "backend/internal/vpn-server"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func New(cfg *config.Config) *App {
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Repositories
	userRepo := user.NewRepository(db)
	serverRepo := vpnserver.NewRepository(db)
	deviceRepo := device.NewRepository(db)

	// Services
	userService := user.NewService(userRepo)
	serverService := vpnserver.NewService(serverRepo)
	deviceService := device.NewService(deviceRepo, serverService)
	authService := auth.NewService(userService, cfg.JWT)

	// Handlers
	userHandler := user.NewHandler(userService)
	serverHandler := vpnserver.NewHandler(serverService)
	deviceHandler := device.NewHandler(deviceService)
	authHandler := auth.NewHandler(authService)

	// Routes
	router := gin.Default()
	router.GET("/health", healthCheck)

	v1 := router.Group("/api/v1")
	auth.RegisterRoutes(v1, authHandler)

	protected := v1.Group("")
	protected.Use(auth.Middleware(authService))

	user.RegisterRoutes(protected, userHandler)
	vpnserver.RegisterRoutes(protected, serverHandler)
	device.RegisterRoutes(protected, deviceHandler)

	return &App{
		Router: router,
		DB:     db,
	}
}
