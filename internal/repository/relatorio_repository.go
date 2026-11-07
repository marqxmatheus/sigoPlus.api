package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"sigo/internal/model"
)

type RelatorioRepository struct {
	db *pgxpool.Pool
}

func NewRelatorioRepository(db *pgxpool.Pool) *RelatorioRepository {
	return &RelatorioRepository{db: db}
}

func (r *RelatorioRepository) Estoque(ctx context.Context, filtro model.FiltroEstoque) ([]model.ItemRelatorioEstoque, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, codigo, descricao, qtd_estoque, estoque_minimo, valor_unitario,
		        qtd_estoque * valor_unitario, qtd_estoque < estoque_minimo
		 FROM pecas
		 WHERE ($1::text = '' OR codigo ILIKE '%' || $1 || '%' OR descricao ILIKE '%' || $1 || '%')
		   AND (NOT $2::boolean OR qtd_estoque < estoque_minimo)
		 ORDER BY descricao ASC, codigo ASC`,
		filtro.Busca, filtro.ApenasAbaixoMinimo,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itens := []model.ItemRelatorioEstoque{}
	for rows.Next() {
		var it model.ItemRelatorioEstoque
		if err := rows.Scan(
			&it.PecaID, &it.Codigo, &it.Descricao, &it.QtdEstoque, &it.EstoqueMinimo,
			&it.ValorUnitario, &it.ValorTotal, &it.AbaixoMinimo,
		); err != nil {
			return nil, err
		}
		itens = append(itens, it)
	}
	return itens, rows.Err()
}
