package service

import (
	"context"
	"errors"
	"regexp"
	"sigo/internal/model"
	"sigo/internal/repository"
	"strings"
)

var (
	ErrNomeObrigatorio   = errors.New("nome é obrigatório")
	ErrDocumentoInvalido = errors.New("CPF/CNPJ deve ter 11 ou 14 dígitos")
)

var naoDigito = regexp.MustCompile(`\D`)

type ClienteService struct {
	repo *repository.ClienteRepository
}

func NewClienteService(repo *repository.ClienteRepository) *ClienteService {
	return &ClienteService{repo: repo}
}

func (s *ClienteService) Listar(ctx context.Context) ([]model.Cliente, error) {
	return s.repo.Listar(ctx)
}

func (s *ClienteService) BuscarPorID(ctx context.Context, id int64) (model.Cliente, error) {
	return s.repo.BuscarPorID(ctx, id)
}

func (s *ClienteService) Criar(ctx context.Context, c *model.Cliente) error {
	if err := normalizarCliente(c); err != nil {
		return err
	}
	return s.repo.Criar(ctx, c)
}

func (s *ClienteService) Atualizar(ctx context.Context, c *model.Cliente) error {
	if err := normalizarCliente(c); err != nil {
		return err
	}
	return s.repo.Atualizar(ctx, c)
}

func (s *ClienteService) Excluir(ctx context.Context, id int64) error {
	return s.repo.Excluir(ctx, id)
}

func normalizarCliente(c *model.Cliente) error {
	c.Nome = strings.TrimSpace(c.Nome)
	c.Email = strings.ToLower(strings.TrimSpace(c.Email))
	c.CPFCNPJ = naoDigito.ReplaceAllString(c.CPFCNPJ, "")

	if c.Nome == "" {
		return ErrNomeObrigatorio
	}
	if len(c.CPFCNPJ) != 11 && len(c.CPFCNPJ) != 14 {
		return ErrDocumentoInvalido
	}
	return nil
}
