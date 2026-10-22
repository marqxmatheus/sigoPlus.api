package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sigo/internal/model"
	"sigo/internal/service"
)

type ClienteHandler struct {
	service *service.ClienteService
}

func NewClienteHandler(s *service.ClienteService) *ClienteHandler {
	return &ClienteHandler{service: s}
}

func (h *ClienteHandler) Listar(c *gin.Context) {
	clientes, err := h.service.Listar(c.Request.Context())
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, clientes)
}

func (h *ClienteHandler) Buscar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	cliente, err := h.service.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, cliente)
}

func (h *ClienteHandler) Criar(c *gin.Context) {
	var cliente model.Cliente
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}
	if err := h.service.Criar(c.Request.Context(), &cliente); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, cliente)
}

func (h *ClienteHandler) Atualizar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	var cliente model.Cliente
	if err := c.ShouldBindJSON(&cliente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dados inválidos: " + err.Error()})
		return
	}
	cliente.ID = id
	if err := h.service.Atualizar(c.Request.Context(), &cliente); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, cliente)
}

func (h *ClienteHandler) Excluir(c *gin.Context) {
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
