package service

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"sigo/internal/model"
	"sigo/internal/repository"
)

// FusoOficina é o horário de Brasília, sem horário de verão desde 2019.
var FusoOficina = time.FixedZone("BRT", -3*60*60)
const formatoData = "2006-01-02"

var ErrPeriodoInvalido = errors.New("a data inicial não pode ser posterior à data final")

type RelatorioService struct {
	repo        *repository.RelatorioRepository
	clienteRepo *repository.ClienteRepository
}

func NewRelatorioService(
	repo *repository.RelatorioRepository,
	clienteRepo *repository.ClienteRepository,
) *RelatorioService {
	return &RelatorioService{repo: repo, clienteRepo: clienteRepo}
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

// OrdensPeriodo recebe datas de calendário; dataFim é inclusiva.
func (s *RelatorioService) OrdensPeriodo(
	ctx context.Context,
	dataInicio, dataFim *time.Time,
) (model.RelatorioOSPeriodo, error) {
	if dataInicio != nil && dataFim != nil && dataInicio.After(*dataFim) {
		return model.RelatorioOSPeriodo{}, ErrPeriodoInvalido
	}

	filtro := model.FiltroOrdens{Inicio: dataInicio}
	if dataFim != nil {
		fimExclusivo := dataFim.AddDate(0, 0, 1)
		filtro.Fim = &fimExclusivo
	}

	itens, err := s.repo.OrdensServico(ctx, filtro)
	if err != nil {
		return model.RelatorioOSPeriodo{}, err
	}

	return model.RelatorioOSPeriodo{
		DataEmissao: time.Now().In(FusoOficina),
		DataInicio:  formatarData(dataInicio),
		DataFim:     formatarData(dataFim),
		Itens:       itens,
		Quantidade:  len(itens),
		ValorTotal:  somarValorOrdens(itens),
	}, nil
}

func (s *RelatorioService) OrdensCliente(
	ctx context.Context,
	clienteID *int64,
) (model.RelatorioOSCliente, error) {
	rel := model.RelatorioOSCliente{DataEmissao: time.Now().In(FusoOficina)}

	if clienteID != nil {
		cliente, err := s.clienteRepo.BuscarPorID(ctx, *clienteID)
		if err != nil {
			return model.RelatorioOSCliente{}, err
		}
		rel.Cliente = &cliente
	}

	itens, err := s.repo.OrdensServico(ctx, model.FiltroOrdens{ClienteID: clienteID})
	if err != nil {
		return model.RelatorioOSCliente{}, err
	}

	rel.Itens = itens
	rel.Quantidade = len(itens)
	rel.ValorTotal = somarValorOrdens(itens)
	return rel, nil
}

func arredondarCentavos(valor float64) float64 {
	return math.Round(valor*100) / 100
}

func ParseData(valor string) (time.Time, error) {
	return time.ParseInLocation(formatoData, valor, FusoOficina)
}

func formatarData(data *time.Time) string {
	if data == nil {
		return ""
	}
	return data.Format(formatoData)
}

func somarValorOrdens(itens []model.ItemRelatorioOS) float64 {
	var total float64
	for _, it := range itens {
		total += it.ValorTotal
	}
	return arredondarCentavos(total)
}
