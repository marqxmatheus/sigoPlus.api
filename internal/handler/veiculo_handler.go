package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sigo/internal/model"
	"sigo/internal/service"
)

type VeiculoHandler struct {
	service *service.VeiculoService
}

func NewVeiculoHandler(s *service.VeiculoService) *VeiculoHandler {
	return &VeiculoHandler{service: s}
}

func (h *VeiculoHandler) Listar(c *gin.Context) {
	veiculos, err := h.service.Listar(c.Request.Context())
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, veiculos)
}

func (h *VeiculoHandler) ListarPorCliente(c *gin.Context) {
	clienteID, ok := lerID(c)
	if !ok {
		return
	}
	veiculos, err := h.service.ListarPorCliente(c.Request.Context(), clienteID)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, veiculos)
}

func (h *VeiculoHandler) Buscar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	veiculo, err := h.service.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, veiculo)
}

func (h *VeiculoHandler) Criar(c *gin.Context) {
	var veiculo model.Veiculo
	if err := c.ShouldBindJSON(&veiculo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	if err := h.service.Criar(c.Request.Context(), &veiculo); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, veiculo)
}

func (h *VeiculoHandler) Atualizar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	var veiculo model.Veiculo
	if err := c.ShouldBindJSON(&veiculo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	veiculo.ID = id
	if err := h.service.Atualizar(c.Request.Context(), &veiculo); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, veiculo)
}

func (h *VeiculoHandler) Excluir(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	if err := h.service.Excluir(c.Request.Context(), id); err != nil {
		responderErro(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}