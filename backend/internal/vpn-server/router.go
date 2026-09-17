package server

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	servers := rg.Group("/server")

	servers.GET("/list", handler.List)
	servers.POST("", handler.Create)
	servers.GET("/:id", handler.Get)
	servers.PUT("/:id", handler.Update)
	servers.DELETE("/:id", handler.Delete)
}
