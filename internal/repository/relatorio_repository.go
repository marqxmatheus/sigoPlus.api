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

func (r *RelatorioRepository) OrdensServico(ctx context.Context, filtro model.FiltroOrdens) ([]model.ItemRelatorioOS, error) {
	rows, err := r.db.Query(ctx,
		`SELECT o.id, o.created_at, o.situacao, o.cliente_id, c.nome,
		        v.placa, v.marca, v.modelo, o.valor_total, o.forma_pagamento, o.data_pagamento
		 FROM ordens_servico o
		 JOIN clientes c ON c.id = o.cliente_id
		 JOIN veiculos v ON v.id = o.veiculo_id
		 WHERE ($1::timestamptz IS NULL OR o.created_at >= $1)
		   AND ($2::timestamptz IS NULL OR o.created_at < $2)
		   AND ($3::bigint IS NULL OR o.cliente_id = $3)
		 ORDER BY o.created_at DESC, o.id DESC`,
		filtro.Inicio, filtro.Fim, filtro.ClienteID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itens := []model.ItemRelatorioOS{}
	for rows.Next() {
		var it model.ItemRelatorioOS
		if err := rows.Scan(
			&it.OrdemServicoID, &it.DataAbertura, &it.Situacao, &it.ClienteID, &it.ClienteNome,
			&it.VeiculoPlaca, &it.VeiculoMarca, &it.VeiculoModelo, &it.ValorTotal,
			&it.FormaPagamento, &it.DataPagamento,
		); err != nil {
			return nil, err
		}
		itens = append(itens, it)
	}
	return itens, rows.Err()
}
