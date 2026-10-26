package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sigo/internal/model"
)

type PecaRepository struct {
	db *pgxpool.Pool
}

func NewPecaRepository(db *pgxpool.Pool) *PecaRepository {
	return &PecaRepository{db: db}
}

const colunasPeca = `id, codigo, descricao, valor_unitario, qtd_estoque, estoque_minimo, created_at, updated_at`

func scanPeca(row pgx.Row) (model.Peca, error) {
	var p model.Peca
	err := row.Scan(&p.ID, &p.Codigo, &p.Descricao, &p.ValorUnitario, &p.QtdEstoque, &p.EstoqueMinimo, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *PecaRepository) Listar(ctx context.Context) ([]model.Peca, error) {
	rows, err := r.db.Query(ctx, `SELECT `+colunasPeca+` FROM pecas ORDER BY descricao ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pecas := []model.Peca{}
	for rows.Next() {
		p, err := scanPeca(rows)
		if err != nil {
			return nil, err
		}
		pecas = append(pecas, p)
	}
	return pecas, rows.Err()
}

func (r *PecaRepository) BuscarPorID(ctx context.Context, id int64) (model.Peca, error) {
	row := r.db.QueryRow(ctx, `SELECT `+colunasPeca+` FROM pecas WHERE id = $1`, id)
	p, err := scanPeca(row)
	return p, traduzirErro(err)
}

func (r *PecaRepository) Criar(ctx context.Context, p *model.Peca) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO pecas (codigo, descricao, valor_unitario, qtd_estoque, estoque_minimo)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		p.Codigo, p.Descricao, p.ValorUnitario, p.QtdEstoque, p.EstoqueMinimo,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return traduzirErro(err)
}

func (r *PecaRepository) Atualizar(ctx context.Context, p *model.Peca) error {
	err := r.db.QueryRow(ctx,
		`UPDATE pecas 
		SET codigo = $1, descricao = $2, valor_unitario = $3, qtd_estoque = $4, estoque_minimo = $5, updated_at = now()
		WHERE id = $6 
		RETURNING created_at, updated_at`,
		p.Codigo, p.Descricao, p.ValorUnitario, p.QtdEstoque, p.EstoqueMinimo, p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
	return traduzirErro(err)
}

func (r *PecaRepository) Excluir(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM pecas WHERE id = $1`, id)
	if err != nil {
		return traduzirErro(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
