package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"sigo/internal/model"
	"sigo/internal/service"
)

type RelatorioHandler struct {
	service *service.RelatorioService
}

func NewRelatorioHandler(s *service.RelatorioService) *RelatorioHandler {
	return &RelatorioHandler{service: s}
}

func (h *RelatorioHandler) Estoque(c *gin.Context) {
	filtro := model.FiltroEstoque{Busca: c.Query("busca")}

	if valor := strings.TrimSpace(c.Query("abaixo_minimo")); valor != "" {
		apenasAbaixo, err := strconv.ParseBool(valor)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "abaixo_minimo deve ser true ou false"})
			return
		}
		filtro.ApenasAbaixoMinimo = apenasAbaixo
	}

	rel, err := h.service.Estoque(c.Request.Context(), filtro)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, rel)
}

func (h *RelatorioHandler) OrdensPeriodo(c *gin.Context) {
	dataInicio, ok := lerData(c, "data_inicio")
	if !ok {
		return
	}
	dataFim, ok := lerData(c, "data_fim")
	if !ok {
		return
	}

	rel, err := h.service.OrdensPeriodo(c.Request.Context(), dataInicio, dataFim)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, rel)
}

func (h *RelatorioHandler) OrdensCliente(c *gin.Context) {
	var clienteID *int64
	if valor := strings.TrimSpace(c.Query("cliente_id")); valor != "" {
		id, err := strconv.ParseInt(valor, 10, 64)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"erro": "cliente_id inválido"})
			return
		}
		clienteID = &id
	}

	rel, err := h.service.OrdensCliente(c.Request.Context(), clienteID)
	if err != nil {
		responderErro(c, err)
		return
	}
	c.JSON(http.StatusOK, rel)
}

func lerData(c *gin.Context, campo string) (*time.Time, bool) {
	valor := strings.TrimSpace(c.Query(campo))
	if valor == "" {
		return nil, true
	}
	data, err := service.ParseData(valor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": campo + " inválida, use o formato AAAA-MM-DD"})
		return nil, false
	}
	return &data, true
}