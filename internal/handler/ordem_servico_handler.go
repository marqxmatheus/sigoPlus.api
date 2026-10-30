package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sigo/internal/model"
	"sigo/internal/service"
)

type OrdemServicoHandler struct {
	service *service.OrdemServicoService
}

func NewOrdemServicoHandler(s *service.OrdemServicoService) *OrdemServicoHandler {
	return &OrdemServicoHandler{service: s}
}

func (h *OrdemServicoHandler) Listar(c *gin.Context) {
	ordens, err := h.service.Listar(c.Request.Context())
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, ordens)
}

func (h *OrdemServicoHandler) Buscar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	ordem, err := h.service.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, ordem)
}

func (h *OrdemServicoHandler) Criar(c *gin.Context) {
	var ordem model.OrdemServico
	if err := c.ShouldBindJSON(&ordem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	if err := h.service.Criar(c.Request.Context(), &ordem); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, ordem)
}

func (h *OrdemServicoHandler) Excluir(c *gin.Context) {
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
