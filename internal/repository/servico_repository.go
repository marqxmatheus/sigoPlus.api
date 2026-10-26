package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sigo/internal/model"
)

type ServicoRepository struct {
	db *pgxpool.Pool
}

func NewServicoRepository(db *pgxpool.Pool) *ServicoRepository {
	return &ServicoRepository{db: db}
}

const colunasServico = `id, descricao, valor_padrao, tempo_estimado_min, created_at, updated_at`

func scanServico(row pgx.Row) (model.Servico, error) {
	var s model.Servico
	err := row.Scan(&s.ID, &s.Descricao, &s.ValorPadrao, &s.TempoEstimadoMin, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *ServicoRepository) Listar(ctx context.Context) ([]model.Servico, error) {
	rows, err := r.db.Query(ctx, `SELECT `+colunasServico+` FROM servicos ORDER BY descricao ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	servicos := []model.Servico{}
	for rows.Next() {
		s, err := scanServico(rows)
		if err != nil {
			return nil, err
		}
		servicos = append(servicos, s)
	}
	return servicos, rows.Err()
}

func (r *ServicoRepository) BuscarPorID(ctx context.Context, id int64) (model.Servico, error) {
	row := r.db.QueryRow(ctx, `SELECT `+colunasServico+` FROM servicos WHERE id = $1`, id)
	s, err := scanServico(row)
	return s, traduzirErro(err)
}

func (r *ServicoRepository) Criar(ctx context.Context, s *model.Servico) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO servicos (descricao, valor_padrao, tempo_estimado_min)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at, updated_at`,
		s.Descricao, s.ValorPadrao, s.TempoEstimadoMin,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	return traduzirErro(err)
}

func (r *ServicoRepository) Atualizar(ctx context.Context, s *model.Servico) error {
	err := r.db.QueryRow(ctx,
		`UPDATE servicos 
		SET descricao = $1, valor_padrao = $2, tempo_estimado_min = $3, updated_at = now()
		WHERE id = $4
		RETURNING created_at, updated_at`,
		s.Descricao, s.ValorPadrao, s.TempoEstimadoMin, s.ID,
	).Scan(&s.CreatedAt, &s.UpdatedAt)
	return traduzirErro(err)
}

func (r *ServicoRepository) Excluir(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM servicos WHERE id = $1`, id)
	if err != nil {
		return traduzirErro(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
