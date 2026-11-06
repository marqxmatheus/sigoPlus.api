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
	ErrItemOSInvalido     = errors.New("item informado não pertence à ordem de serviço")
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

func (s *OrdemServicoService) Atualizar(
	ctx context.Context,
	id int64,
	atualizacao *model.OrdemServicoAtualizacao,
) (model.OrdemServico, error) {
	atual, err := s.repo.BuscarPorID(ctx, id)
	if err != nil {
		return model.OrdemServico{}, err
	}

	if len(atualizacao.Servicos) == 0 && len(atualizacao.Pecas) == 0 {
		return model.OrdemServico{}, ErrOSSemItens
	}

	if err := s.prepararServicosAtualizacao(ctx, atual.Servicos, atualizacao.Servicos); err != nil {
		return model.OrdemServico{}, err
	}
	if err := s.prepararPecasAtualizacao(ctx, atual.Pecas, atualizacao.Pecas); err != nil {
		return model.OrdemServico{}, err
	}

	atual.Observacao = strings.TrimSpace(atualizacao.Observacao)
	atual.Servicos = atualizacao.Servicos
	atual.Pecas = atualizacao.Pecas

	if err := s.repo.Atualizar(ctx, &atual); err != nil {
		return model.OrdemServico{}, err
	}

	return s.repo.BuscarPorID(ctx, id)
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

func (s *OrdemServicoService) prepararServicosAtualizacao(
	ctx context.Context,
	atuais []model.OrdemServicoItemServico,
	novos []model.OrdemServicoItemServico,
) error {
	usados := make(map[int64]bool, len(atuais))

	for i := range novos {
		servico, err := s.servicoRepo.BuscarPorID(ctx, novos[i].ServicoID)
		if errors.Is(err, repository.ErrNaoEncontrado) {
			return fmt.Errorf("%w (id %d)", ErrServicoInexistente, novos[i].ServicoID)
		}
		if err != nil {
			return err
		}

		indiceAtual := encontrarServicoAtual(atuais, novos[i], usados)
		if novos[i].ID != 0 && indiceAtual == -1 {
			return fmt.Errorf("%w (serviço id %d)", ErrItemOSInvalido, novos[i].ID)
		}

		if indiceAtual >= 0 {
			itemAtual := atuais[indiceAtual]
			usados[itemAtual.ID] = true
			novos[i].ID = itemAtual.ID
			novos[i].Valor = itemAtual.Valor
		} else {
			novos[i].ID = 0
			novos[i].Valor = servico.ValorPadrao
		}
	}

	return nil
}

func (s *OrdemServicoService) prepararPecasAtualizacao(
	ctx context.Context,
	atuais []model.OrdemServicoItemPeca,
	novas []model.OrdemServicoItemPeca,
) error {
	usadas := make(map[int64]bool, len(atuais))

	for i := range novas {
		peca, err := s.pecaRepo.BuscarPorID(ctx, novas[i].PecaID)
		if errors.Is(err, repository.ErrNaoEncontrado) {
			return fmt.Errorf("%w (id %d)", ErrPecaInexistente, novas[i].PecaID)
		}
		if err != nil {
			return err
		}

		indiceAtual := encontrarPecaAtual(atuais, novas[i], usadas)
		if novas[i].ID != 0 && indiceAtual == -1 {
			return fmt.Errorf("%w (peça id %d)", ErrItemOSInvalido, novas[i].ID)
		}

		if indiceAtual >= 0 {
			itemAtual := atuais[indiceAtual]
			usadas[itemAtual.ID] = true
			novas[i].ID = itemAtual.ID
			novas[i].ValorUnitario = itemAtual.ValorUnitario
		} else {
			novas[i].ID = 0
			novas[i].ValorUnitario = peca.ValorUnitario
		}
	}

	return nil
}

func encontrarServicoAtual(
	atuais []model.OrdemServicoItemServico,
	novo model.OrdemServicoItemServico,
	usados map[int64]bool,
) int {
	for i := range atuais {
		if usados[atuais[i].ID] {
			continue
		}
		if novo.ID > 0 && atuais[i].ID == novo.ID && atuais[i].ServicoID == novo.ServicoID {
			return i
		}
	}

	if novo.ID != 0 {
		return -1
	}

	for i := range atuais {
		if !usados[atuais[i].ID] && atuais[i].ServicoID == novo.ServicoID {
			return i
		}
	}

	return -1
}

func encontrarPecaAtual(
	atuais []model.OrdemServicoItemPeca,
	nova model.OrdemServicoItemPeca,
	usadas map[int64]bool,
) int {
	for i := range atuais {
		if usadas[atuais[i].ID] {
			continue
		}
		if nova.ID > 0 && atuais[i].ID == nova.ID && atuais[i].PecaID == nova.PecaID {
			return i
		}
	}

	if nova.ID != 0 {
		return -1
	}

	for i := range atuais {
		if !usadas[atuais[i].ID] && atuais[i].PecaID == nova.PecaID {
			return i
		}
	}

	return -1
}
