package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sigo/internal/model"
	"sigo/internal/service"
)

type ServicoHandler struct {
	service *service.ServicoService
}

func NewServicoHandler(s *service.ServicoService) *ServicoHandler {
	return &ServicoHandler{service: s}
}

func (h *ServicoHandler) Listar(c *gin.Context) {
	servicos, err := h.service.Listar(c.Request.Context())
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, servicos)
}

func (h *ServicoHandler) Buscar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	servico, err := h.service.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, servico)
}

func (h *ServicoHandler) Criar(c *gin.Context) {
	var servico model.Servico
	if err := c.ShouldBindJSON(&servico); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	if err := h.service.Criar(c.Request.Context(), &servico); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, servico)
}

func (h *ServicoHandler) Atualizar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	var servico model.Servico
	if err := c.ShouldBindJSON(&servico); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	servico.ID = id
	if err := h.service.Atualizar(c.Request.Context(), &servico); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, servico)
}

func (h *ServicoHandler) Excluir(c *gin.Context) {
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
