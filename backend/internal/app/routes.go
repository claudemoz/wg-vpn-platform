package app

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userService, cfg.JWT)
	authHandler := auth.NewHandler(authService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
		}

		protected := v1.Group("")
		protected.Use(auth.Middleware(authService))
		{
			protected.GET("/user/me", userHandler.GetMe)
			protected.PUT("/user/me", userHandler.UpdateMe)

			protected.GET("/user/list", userHandler.List)
			protected.GET("/user/:id", userHandler.Get)
			protected.PUT("/user/:id", userHandler.Update)
			protected.DELETE("/user/:id", userHandler.Delete)
		}
	}
}
