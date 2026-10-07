package model

import "time"

type Servico struct {
	ID               int64     `json:"id"`
	Descricao        string    `json:"descricao" binding:"required"`
	ValorPadrao      float64   `json:"valor_padrao" binding:"gte=0"`
	TempoEstimadoMin int       `json:"tempo_estimado_min" binding:"gte=0"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
