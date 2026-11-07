package handler

import (
	"net/http"
	"strconv"
	"strings"

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
