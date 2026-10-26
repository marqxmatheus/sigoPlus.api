package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sigo/internal/model"
	"sigo/internal/service"
)

type PecaHandler struct {
	service *service.PecaService
}

func NewPecaHandler(s *service.PecaService) *PecaHandler {
	return &PecaHandler{service: s}
}

func (h *PecaHandler) Listar(c *gin.Context) {
	pecas, err := h.service.Listar(c.Request.Context())
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, pecas)
}

func (h *PecaHandler) Buscar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	peca, err := h.service.BuscarPorID(c.Request.Context(), id)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, peca)
}

func (h *PecaHandler) Criar(c *gin.Context) {
	var peca model.Peca
	if err := c.ShouldBindJSON(&peca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	if err := h.service.Criar(c.Request.Context(), &peca); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusCreated, peca)
}

func (h *PecaHandler) Atualizar(c *gin.Context) {
	id, ok := lerID(c)
	if !ok {
		return
	}
	var peca model.Peca
	if err := c.ShouldBindJSON(&peca); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "dados inválidos: " + err.Error()})
		return
	}
	peca.ID = id
	if err := h.service.Atualizar(c.Request.Context(), &peca); err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, peca)
}

func (h *PecaHandler) Excluir(c *gin.Context) {
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
