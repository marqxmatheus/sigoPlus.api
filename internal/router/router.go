package router

import (
	"github.com/gin-gonic/gin"

	"sigo/internal/config"
	"sigo/internal/handler"
)

type Handlers struct {
	Cliente      *handler.ClienteHandler
	Veiculo      *handler.VeiculoHandler
	Servico      *handler.ServicoHandler
	Peca         *handler.PecaHandler
	OrdemServico *handler.OrdemServicoHandler
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
			clientes.GET("/:id/veiculos", h.Veiculo.ListarPorCliente)
			clientes.POST("", h.Cliente.Criar)
			clientes.PUT("/:id", h.Cliente.Atualizar)
			clientes.DELETE("/:id", h.Cliente.Excluir)
		}

		veiculos := api.Group("/veiculos")
		{
			veiculos.GET("", h.Veiculo.Listar)
			veiculos.GET("/:id", h.Veiculo.Buscar)
			veiculos.POST("", h.Veiculo.Criar)
			veiculos.PUT("/:id", h.Veiculo.Atualizar)
			veiculos.DELETE("/:id", h.Veiculo.Excluir)
		}

		servicos := api.Group("/servicos")
		{
			servicos.GET("", h.Servico.Listar)
			servicos.GET("/:id", h.Servico.Buscar)
			servicos.POST("", h.Servico.Criar)
			servicos.PUT("/:id", h.Servico.Atualizar)
			servicos.DELETE("/:id", h.Servico.Excluir)
		}

		pecas := api.Group("/pecas")
		{
			pecas.GET("", h.Peca.Listar)
			pecas.GET("/:id", h.Peca.Buscar)
			pecas.POST("", h.Peca.Criar)
			pecas.PUT("/:id", h.Peca.Atualizar)
			pecas.DELETE("/:id", h.Peca.Excluir)
		}

		ordens := api.Group("/ordens-servico")
		{
			ordens.GET("", h.OrdemServico.Listar)
			ordens.GET("/:id", h.OrdemServico.Buscar)
			ordens.POST("", h.OrdemServico.Criar)
			ordens.DELETE("/:id", h.OrdemServico.Excluir)
		}
	}
	return r
}
