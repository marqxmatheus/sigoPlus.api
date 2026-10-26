package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sigo/internal/model"
)

type VeiculoRepository struct {
	db *pgxpool.Pool
}

func NewVeiculoRepository(db *pgxpool.Pool) *VeiculoRepository {
	return &VeiculoRepository{db: db}
}

const colunasVeiculo = `id, cliente_id, placa, marca, modelo, ano, created_at, updated_at`

func scanVeiculo(row pgx.Row) (model.Veiculo, error) {
	var v model.Veiculo
	err := row.Scan(&v.ID, &v.ClienteID, &v.Placa, &v.Marca, &v.Modelo, &v.Ano, &v.CreatedAt, &v.UpdatedAt)
	return v, err
}

func (r *VeiculoRepository) listar(ctx context.Context, sql string, args ...any) ([]model.Veiculo, error) {
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	veiculos := []model.Veiculo{}
	for rows.Next() {
		v, err := scanVeiculo(rows)
		if err != nil {
			return nil, err
		}
		veiculos = append(veiculos, v)
	}
	return veiculos, rows.Err()
}

func (r *VeiculoRepository) Listar(ctx context.Context) ([]model.Veiculo, error) {
	return r.listar(ctx, `SELECT `+colunasVeiculo+` FROM veiculos ORDER BY placa`)
}

func (r *VeiculoRepository) ListarPorCliente(ctx context.Context, clienteID int64) ([]model.Veiculo, error) {
	return r.listar(ctx, `SELECT `+colunasVeiculo+` FROM veiculos WHERE cliente_id = $1 ORDER BY placa`, clienteID)
}

func (r *VeiculoRepository) BuscarPorID(ctx context.Context, id int64) (model.Veiculo, error) {
	row := r.db.QueryRow(ctx, `SELECT `+colunasVeiculo+` FROM veiculos WHERE id = $1`, id)
	v, err := scanVeiculo(row)
	return v, traduzirErro(err)
}

func (r *VeiculoRepository) Criar(ctx context.Context, v *model.Veiculo) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO veiculos (cliente_id, placa, marca, modelo, ano)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		v.ClienteID, v.Placa, v.Marca, v.Modelo, v.Ano,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
	return traduzirErro(err)
}

func (r *VeiculoRepository) Atualizar(ctx context.Context, v *model.Veiculo) error {
	err := r.db.QueryRow(ctx,
		`UPDATE veiculos
		 SET cliente_id = $1, placa = $2, marca = $3, modelo = $4, ano = $5, updated_at = now()
		 WHERE id = $6
		 RETURNING created_at, updated_at`,
		v.ClienteID, v.Placa, v.Marca, v.Modelo, v.Ano, v.ID,
	).Scan(&v.CreatedAt, &v.UpdatedAt)
	return traduzirErro(err)
}

func (r *VeiculoRepository) Excluir(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM veiculos WHERE id = $1`, id)
	if err != nil {
		return traduzirErro(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}
