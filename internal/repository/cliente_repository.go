package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sigo/internal/model"
)

type ClienteRepository struct {
	db *pgxpool.Pool
}

func NewClienteRepository(db *pgxpool.Pool) *ClienteRepository {
	return &ClienteRepository{db: db}
}

const colunasCliente = `id, nome, cpf_cnpj, telefone, email, endereco, created_at, updated_at`

func scanCliente(row pgx.Row) (model.Cliente, error) {
	var c model.Cliente
	err := row.Scan(&c.ID, &c.Nome, &c.CPFCNPJ, &c.Telefone, &c.Email, &c.Endereco, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (r *ClienteRepository) Listar(ctx context.Context) ([]model.Cliente, error) {
    rows, err := r.db.Query(ctx, `SELECT `+colunasCliente+` FROM clientes ORDER BY nome ASC`)
    if err != nil {
		return nil, err
	}
	defer rows.Close()

	clientes := []model.Cliente{}
	for rows.Next() {
		c, err := scanCliente(rows)
		if err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, rows.Err()
}

func (r *ClienteRepository) BuscarPorID(ctx context.Context, id int64) (model.Cliente, error) {
	row := r.db.QueryRow(ctx, `SELECT `+colunasCliente+` FROM clientes WHERE id = $1`, id)
	c, err := scanCliente(row)
	return c, traduzirErro(err)
}

func (r *ClienteRepository) Criar(ctx context.Context, c *model.Cliente) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO clientes (nome, cpf_cnpj, telefone, email, endereco)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		c.Nome, c.CPFCNPJ, c.Telefone, c.Email, c.Endereco,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	return traduzirErro(err)
}

func (r *ClienteRepository) Atualizar(ctx context.Context, c *model.Cliente) error {
	err := r.db.QueryRow(ctx,
		`UPDATE clientes 
		SET nome = $1, cpf_cnpj = $2, telefone = $3, email = $4, endereco = $5, updated_at = now()
		WHERE id = $6 
		RETURNING created_at, updated_at`,
		c.Nome, c.CPFCNPJ, c.Telefone, c.Email, c.Endereco, c.ID,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
	return traduzirErro(err)
}

func (r *ClienteRepository) Excluir(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM clientes WHERE id = $1`, id)
	if err != nil{
		return traduzirErro(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNaoEncontrado
	}
	return nil
}


