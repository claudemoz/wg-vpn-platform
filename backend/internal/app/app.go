package app

import (
	"backend/internal/config"
	"backend/internal/database"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB    *gorm.DB
}

func New(cfg *config.Config) *App {
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	RegisterRoutes(router, db, cfg)

	return &App{
		Router: router,
		DB:     db,
	}
}
