package service

import (
	"context"
	"errors"
	"sigo/internal/model"
	"sigo/internal/repository"
	"strings"
)

var ErrDescricaoObrigatoria = errors.New("descrição é obrigatória")

type ServicoService struct {
	repo *repository.ServicoRepository
}

func NewServicoService(repo *repository.ServicoRepository) *ServicoService {
	return &ServicoService{repo: repo}
}

func (s *ServicoService) Listar(ctx context.Context) ([]model.Servico, error) {
	return s.repo.Listar(ctx)
}

func (s *ServicoService) BuscarPorID(ctx context.Context, id int64) (model.Servico, error) {
	return s.repo.BuscarPorID(ctx, id)
}

func (s *ServicoService) Criar(ctx context.Context, c *model.Servico) error {
	if err := normalizarServico(c); err != nil {
		return err
	}
	return s.repo.Criar(ctx, c)
}

func (s *ServicoService) Atualizar(ctx context.Context, c *model.Servico) error {
	if err := normalizarServico(c); err != nil {
		return err
	}
	return s.repo.Atualizar(ctx, c)
}

func (s *ServicoService) Excluir(ctx context.Context, id int64) error {
	return s.repo.Excluir(ctx, id)
}

func normalizarServico(s *model.Servico) error {
	s.Descricao = strings.TrimSpace(s.Descricao)
	if s.Descricao == "" {
		return ErrDescricaoObrigatoria
	}
	return nil
}
