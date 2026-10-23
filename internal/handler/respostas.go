package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"sigo/internal/repository"
	"sigo/internal/service"
)

func lerID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
		return 0, false
	}
	return id, true
}

func responderErro(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNaoEncontrado):
		c.JSON(http.StatusNotFound, gin.H{"erro": err.Error()})
	case errors.Is(err, repository.ErrDuplicado), errors.Is(err, repository.ErrEmUso):
		c.JSON(http.StatusConflict, gin.H{"erro": err.Error()})
	case errors.Is(err, service.ErrNomeObrigatorio), errors.Is(err, service.ErrDocumentoInvalido):
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
	default:
		log.Printf("Erro interno: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro interno do servidor"})
	}
}
