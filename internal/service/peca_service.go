package service

import (
	"context"
	"sigo/internal/model"
	"sigo/internal/repository"
	"strings"
)

type PecaService struct {
	repo *repository.PecaRepository
}

func NewPecaService(repo *repository.PecaRepository) *PecaService {
	return &PecaService{repo: repo}
}

func (s *PecaService) Listar(ctx context.Context) ([]model.Peca, error) {
	return s.repo.Listar(ctx)
}

func (s *PecaService) BuscarPorID(ctx context.Context, id int64) (model.Peca, error) {
	return s.repo.BuscarPorID(ctx, id)
}

func (s *PecaService) Criar(ctx context.Context, c *model.Peca) error {
	if err := normalizarPeca(c); err != nil {
		return err
	}
	return s.repo.Criar(ctx, c)
}

func (s *PecaService) Atualizar(ctx context.Context, c *model.Peca) error {
	if err := normalizarPeca(c); err != nil {
		return err
	}
	return s.repo.Atualizar(ctx, c)
}

func (s *PecaService) Excluir(ctx context.Context, id int64) error {
	return s.repo.Excluir(ctx, id)
}

func normalizarPeca(p *model.Peca) error {
	p.Codigo = strings.ToUpper(strings.TrimSpace(p.Codigo))
	p.Descricao = strings.TrimSpace(p.Descricao)
	if p.Descricao == "" {
		return ErrDescricaoObrigatoria
	}
	return nil
}
