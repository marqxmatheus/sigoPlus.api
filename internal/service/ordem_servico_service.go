package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"sigo/internal/model"
	"sigo/internal/repository"
)

var (
	ErrOSSemItens         = errors.New("a ordem de serviço precisa ter pelo menos um serviço ou peça")
	ErrVeiculoInvalido    = errors.New("veículo não existe ou não pertence ao cliente")
	ErrServicoInexistente = errors.New("serviço informado não existe")
	ErrPecaInexistente    = errors.New("peça informada não existe")
	ErrOSNaoExcluivel     = errors.New("apenas ordens de serviço em ORCADO podem ser excluídas")
)

type OrdemServicoService struct {
	repo        *repository.OrdemServicoRepository
	clienteRepo *repository.ClienteRepository
	veiculoRepo *repository.VeiculoRepository
	servicoRepo *repository.ServicoRepository
	pecaRepo    *repository.PecaRepository
}

func NewOrdemServicoService(
	repo *repository.OrdemServicoRepository,
	clienteRepo *repository.ClienteRepository,
	veiculoRepo *repository.VeiculoRepository,
	servicoRepo *repository.ServicoRepository,
	pecaRepo *repository.PecaRepository,
) *OrdemServicoService {
	return &OrdemServicoService{
		repo:        repo,
		clienteRepo: clienteRepo,
		veiculoRepo: veiculoRepo,
		servicoRepo: servicoRepo,
		pecaRepo:    pecaRepo,
	}
}

func (s *OrdemServicoService) Listar(ctx context.Context) ([]model.OrdemServico, error) {
	return s.repo.Listar(ctx)
}

func (s *OrdemServicoService) BuscarPorID(ctx context.Context, id int64) (model.OrdemServico, error) {
	return s.repo.BuscarPorID(ctx, id)
}

func (s *OrdemServicoService) Criar(ctx context.Context, ordem *model.OrdemServico) error {
	if len(ordem.Servicos) == 0 && len(ordem.Pecas) == 0 {
		return ErrOSSemItens
	}
	if err := s.validarClienteEVeiculo(ctx, ordem); err != nil {
		return err
	}
	if err := s.preencherServicos(ctx, ordem.Servicos); err != nil {
		return err
	}
	if err := s.preencherPecas(ctx, ordem.Pecas); err != nil {
		return err
	}

	ordem.Situacao = model.SituacaoOrcado
	ordem.Observacao = strings.TrimSpace(ordem.Observacao)
	ordem.FormaPagamento = nil
	ordem.DataPagamento = nil

	if err := s.repo.Criar(ctx, ordem); err != nil {
		return err
	}

	completa, err := s.repo.BuscarPorID(ctx, ordem.ID)
	if err != nil {
		return err
	}
	*ordem = completa
	return nil
}

func (s *OrdemServicoService) Excluir(ctx context.Context, id int64) error {
	ordem, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if ordem.Situacao != model.SituacaoOrcado {
		return ErrOSNaoExcluivel
	}
	return s.repo.Excluir(ctx, id)
}

func (s *OrdemServicoService) validarClienteEVeiculo(ctx context.Context, ordem *model.OrdemServico) error {
	if _, err := s.clienteRepo.BuscarPorID(ctx, ordem.ClienteID); err != nil {
		if errors.Is(err, repository.ErrNaoEncontrado) {
			return ErrClienteInexistente
		}
		return err
	}

	veiculo, err := s.veiculoRepo.BuscarPorID(ctx, ordem.VeiculoID)
	if errors.Is(err, repository.ErrNaoEncontrado) || (err == nil && veiculo.ClienteID != ordem.ClienteID) {
		return ErrVeiculoInvalido
	}
	return err
}

func (s *OrdemServicoService) preencherServicos(ctx context.Context, itens []model.OrdemServicoItemServico) error {
	for i := range itens {
		servico, err := s.servicoRepo.BuscarPorID(ctx, itens[i].ServicoID)
		if errors.Is(err, repository.ErrNaoEncontrado) {
			return fmt.Errorf("%w (id %d)", ErrServicoInexistente, itens[i].ServicoID)
		}
		if err != nil {
			return err
		}
		itens[i].Valor = servico.ValorPadrao
	}
	return nil
}

func (s *OrdemServicoService) preencherPecas(ctx context.Context, itens []model.OrdemServicoItemPeca) error {
	for i := range itens {
		peca, err := s.pecaRepo.BuscarPorID(ctx, itens[i].PecaID)
		if errors.Is(err, repository.ErrNaoEncontrado) {
			return fmt.Errorf("%w (id %d)", ErrPecaInexistente, itens[i].PecaID)
		}
		if err != nil {
			return err
		}
		itens[i].ValorUnitario = peca.ValorUnitario
	}
	return nil
}
