package model

import "time"

type Peca struct {
	ID            int64     `json:"id"`
	Codigo        string    `json:"codigo" binding:"required"`
	Descricao     string    `json:"descricao" binding:"required"`
	ValorUnitario float64   `json:"valor_unitario" binding:"gte=0"`
	QtdEstoque    int       `json:"qtd_estoque" binding:"gte=0"`
	EstoqueMinimo int       `json:"estoque_minimo" binding:"gte=0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (p Peca) AbaixoDoMinimo() bool {
	return p.QtdEstoque < p.EstoqueMinimo
}
