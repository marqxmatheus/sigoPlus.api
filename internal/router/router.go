package router

import (
	"github.com/gin-gonic/gin"

	"sigo/internal/config"
	"sigo/internal/handler"
)

type Handlers struct {
	Cliente *handler.ClienteHandler
}

func New(cfg *config.Config, h Handlers) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.Health)

		clientes := api.Group("/clientes")
		{
			clientes.GET("", h.Cliente.Listar)
			clientes.GET("/:id", h.Cliente.Buscar)
			clientes.POST("", h.Cliente.Criar)
			clientes.PUT("/:id", h.Cliente.Atualizar)
			clientes.DELETE("/:id", h.Cliente.Excluir)
		}
	}

	return r
}
