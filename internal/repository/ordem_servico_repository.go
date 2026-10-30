package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sigo/internal/model"
)

type OrdemServicoRepository struct {
	db *pgxpool.Pool
}

func NewOrdemServicoRepository(db *pgxpool.Pool) *OrdemServicoRepository {
	return &OrdemServicoRepository{db: db}
}

const colunasOS = `o.id, o.cliente_id, c.nome, o.veiculo_id, v.placa, o.situacao, o.valor_total,
	o.observacao, o.forma_pagamento, o.data_pagamento, o.created_at, o.updated_at`

const fromOS = ` FROM ordens_servico o
	JOIN clientes c ON c.id = o.cliente_id
	JOIN veiculos v ON v.id = o.veiculo_id`

const sqlRecalcularTotal = `
	UPDATE ordens_servico
	SET valor_total =
		  (SELECT COALESCE(SUM(valor), 0) FROM os_servicos WHERE ordem_servico_id = $1)
		+ (SELECT COALESCE(SUM(quantidade * valor_unitario), 0) FROM os_pecas WHERE ordem_servico_id = $1)
	WHERE id = $1
	RETURNING valor_total`

func scanOS(row pgx.Row) (model.OrdemServico, error) {
	var o model.OrdemServico
	err := row.Scan(&o.ID, &o.ClienteID, &o.ClienteNome, &o.VeiculoID, &o.VeiculoPlaca, &o.Situacao,
		&o.ValorTotal, &o.Observacao, &o.FormaPagamento, &o.DataPagamento, &o.CreatedAt, &o.UpdatedAt)
	return o, err
}

func (r *OrdemServicoRepository) Listar(ctx context.Context) ([]model.OrdemServico, error) {
	rows, err := r.db.Query(ctx, `SELECT `+colunasOS+fromOS+` ORDER BY o.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordens := []model.OrdemServico{}
	for rows.Next() {
		o, err := scanOS(rows)
		if err != nil {
			return nil, err
		}
		ordens = append(ordens, o)
	}
	return ordens, rows.Err()
}

func (r *OrdemServicoRepository) BuscarPorID(ctx context.Context, id int64) (model.OrdemServico, error) {
	o, err := scanOS(r.db.QueryRow(ctx, `SELECT `+colunasOS+fromOS+` WHERE o.id = $1`, id))
	if err != nil {
		return o, traduzirErro(err)
	}

	o.Servicos, err = r.listarServicos(ctx, id)
	if err != nil {
		return o, err
	}
	o.Pecas, err = r.listarPecas(ctx, id)
	return o, err
}

func (r *OrdemServicoRepository) listarServicos(ctx context.Context, ordemID int64) ([]model.OrdemServicoItemServico, error) {
	rows, err := r.db.Query(ctx,
		`SELECT i.id, i.ordem_servico_id, i.servico_id, s.descricao, i.valor
		 FROM os_servicos i
		 JOIN servicos s ON s.id = i.servico_id
		 WHERE i.ordem_servico_id = $1
		 ORDER BY i.id`, ordemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itens := []model.OrdemServicoItemServico{}
	for rows.Next() {
		var it model.OrdemServicoItemServico
		if err := rows.Scan(&it.ID, &it.OrdemServicoID, &it.ServicoID, &it.Descricao, &it.Valor); err != nil {
			return nil, err
		}
		itens = append(itens, it)
	}
	return itens, rows.Err()
}

func (r *OrdemServicoRepository) listarPecas(ctx context.Context, ordemID int64) ([]model.OrdemServicoItemPeca, error) {
	rows, err := r.db.Query(ctx,
		`SELECT i.id, i.ordem_servico_id, i.peca_id, p.descricao, i.quantidade, i.valor_unitario
		 FROM os_pecas i
		 JOIN pecas p ON p.id = i.peca_id
		 WHERE i.ordem_servico_id = $1
		 ORDER BY i.id`, ordemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itens := []model.OrdemServicoItemPeca{}
	for rows.Next() {
		var it model.OrdemServicoItemPeca
		if err := rows.Scan(&it.ID, &it.OrdemServicoID, &it.PecaID, &it.Descricao, &it.Quantidade, &it.ValorUnitario); err != nil {
			return nil, err
		}
		itens = append(itens, it)
	}
	return itens, rows.Err()
}

func (r *OrdemServicoRepository) Criar(ctx context.Context, ordem *model.OrdemServico) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO ordens_servico (cliente_id, veiculo_id, situacao, observacao)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, created_at, updated_at`,
		ordem.ClienteID, ordem.VeiculoID, ordem.Situacao, ordem.Observacao,
	).Scan(&ordem.ID, &ordem.CreatedAt, &ordem.UpdatedAt)
	if err != nil {
		return traduzirErro(err)
	}

	for i := range ordem.Servicos {
		item := &ordem.Servicos[i]
		item.OrdemServicoID = ordem.ID
		err = tx.QueryRow(ctx,
			`INSERT INTO os_servicos (ordem_servico_id, servico_id, valor)
			 VALUES ($1, $2, $3)
			 RETURNING id`,
			ordem.ID, item.ServicoID, item.Valor,
		).Scan(&item.ID)
		if err != nil {
			return traduzirErro(err)
		}
	}

	for i := range ordem.Pecas {
		item := &ordem.Pecas[i]
		item.OrdemServicoID = ordem.ID
		err = tx.QueryRow(ctx,
			`INSERT INTO os_pecas (ordem_servico_id, peca_id, quantidade, valor_unitario)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			ordem.ID, item.PecaID, item.Quantidade, item.ValorUnitario,
		).Scan(&item.ID)
		if err != nil {
			return traduzirErro(err)
		}
	}

	if err := tx.QueryRow(ctx, sqlRecalcularTotal, ordem.ID).Scan(&ordem.ValorTotal); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrdemServicoRepository) Excluir(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM ordens_servico WHERE id = $1`, id)
	if err != nil {
		return traduzirErro(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
