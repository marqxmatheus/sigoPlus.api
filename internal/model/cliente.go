package model

import "time"

type Cliente struct {
	ID        int64     `json:"id"`
	Nome      string    `json:"nome" binding:"required"`
	CPFCNPJ   string    `json:"cpf_cnpj" binding:"required"`
	Telefone  string    `json:"telefone"`
	Email     string    `json:"email" binding:"omitempty,email"`
	Endereco  string    `json:"endereco"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
