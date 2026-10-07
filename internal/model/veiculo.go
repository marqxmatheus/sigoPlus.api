package model

import "time"

type Veiculo struct {
	ID        int64     `json:"id"`
	ClienteID int64     `json:"cliente_id" binding:"required"`
	Placa     string    `json:"placa" binding:"required"`
	Marca     string    `json:"marca"`
	Modelo    string    `json:"modelo"`
	Ano       int       `json:"ano"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
