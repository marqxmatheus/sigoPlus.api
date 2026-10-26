package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"sigo/internal/model"
	"sigo/internal/repository"
)

var (
	ErrPlacaInvalida      = errors.New("placa deve estar no formato ABC1234 ou ABC1D23")
	ErrClienteInexistente = errors.New("cliente informado não existe")
)

var (
	formatoPlaca    = regexp.MustCompile(`^[A-Z]{3}[0-9][A-Z0-9][0-9]{2}$`)
	naoAlfanumerico = regexp.MustCompile(`[^A-Z0-9]`)
)

type VeiculoService struct {
	repo        *repository.VeiculoRepository
	clienteRepo *repository.ClienteRepository
}

func NewVeiculoService(repo *repository.VeiculoRepository, clienteRepo *repository.ClienteRepository) *VeiculoService {
	return &VeiculoService{repo: repo, clienteRepo: clienteRepo}
}

func (s *VeiculoService) Listar(ctx context.Context) ([]model.Veiculo, error) {
	return s.repo.Listar(ctx)
}

func (s *VeiculoService) ListarPorCliente(ctx context.Context, clienteID int64) ([]model.Veiculo, error) {
	if _, err := s.clienteRepo.BuscarPorID(ctx, clienteID); err != nil {
		return nil, err
	}
	return s.repo.ListarPorCliente(ctx, clienteID)
}

func (s *VeiculoService) BuscarPorID(ctx context.Context, id int64) (model.Veiculo, error) {
	return s.repo.BuscarPorID(ctx, id)
}

func (s *VeiculoService) Criar(ctx context.Context, v *model.Veiculo) error {
	if err := s.validar(ctx, v); err != nil {
		return err
	}
	return s.repo.Criar(ctx, v)
}

func (s *VeiculoService) Atualizar(ctx context.Context, v *model.Veiculo) error {
	if err := s.validar(ctx, v); err != nil {
		return err
	}
	return s.repo.Atualizar(ctx, v)
}

func (s *VeiculoService) Excluir(ctx context.Context, id int64) error {
	return s.repo.Excluir(ctx, id)
}

func (s *VeiculoService) validar(ctx context.Context, v *model.Veiculo) error {
	v.Placa = naoAlfanumerico.ReplaceAllString(strings.ToUpper(v.Placa), "")
	v.Marca = strings.TrimSpace(v.Marca)
	v.Modelo = strings.TrimSpace(v.Modelo)

	if !formatoPlaca.MatchString(v.Placa) {
		return ErrPlacaInvalida
	}

	_, err := s.clienteRepo.BuscarPorID(ctx, v.ClienteID)
	if errors.Is(err, repository.ErrNaoEncontrado) {
		return ErrClienteInexistente
	}
	return err
}
