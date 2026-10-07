package router

import (
	"github.com/gin-gonic/gin"

	"sigo/internal/config"
	"sigo/internal/handler"
)

func New(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.Health)
	}

	return r
}
