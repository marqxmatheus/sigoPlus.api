package repository

import (
	"errors"
	
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNaoEncontrado = errors.New("registro não encontrado")
	ErrDuplicado = errors.New("Já existe um registro com esses dados")
	ErrEmUso = errors.New("registro possui vínculos e não pode ser excluído")
)

func traduzirErro(err error) error {
	if err == nil{
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows){
		return ErrNaoEncontrado
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr){
		switch pgErr.Code {
			case "23505":
				return ErrDuplicado
			case "23503":
				return ErrEmUso
		}
	}
	return err
}