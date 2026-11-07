package service

import (
	"context"
	"math"
	"strings"
	"time"

	"sigo/internal/model"
	"sigo/internal/repository"
)

// FusoOficina é o horário de Brasília, sem horário de verão desde 2019.
var FusoOficina = time.FixedZone("BRT", -3*60*60)

type RelatorioService struct {
	repo *repository.RelatorioRepository
}

func NewRelatorioService(repo *repository.RelatorioRepository) *RelatorioService {
	return &RelatorioService{repo: repo}
}

func (s *RelatorioService) Estoque(ctx context.Context, filtro model.FiltroEstoque) (model.RelatorioEstoque, error) {
	filtro.Busca = strings.TrimSpace(filtro.Busca)

	itens, err := s.repo.Estoque(ctx, filtro)
	if err != nil {
		return model.RelatorioEstoque{}, err
	}

	rel := model.RelatorioEstoque{
		DataEmissao: time.Now().In(FusoOficina),
		Itens:       itens,
		TotalPecas:  len(itens),
	}
	for _, it := range itens {
		rel.QuantidadeTotal += it.QtdEstoque
		rel.ValorTotal += it.ValorTotal
		if it.AbaixoMinimo {
			rel.TotalAbaixoMinimo++
		}
	}
	rel.ValorTotal = arredondarCentavos(rel.ValorTotal)
	return rel, nil
}

func arredondarCentavos(valor float64) float64 {
	return math.Round(valor*100) / 100
}
