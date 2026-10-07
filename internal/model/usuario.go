package model

import "time"

type Perfil string

const (
	PerfilGestor    Perfil = "GESTOR"
	PerfilAtendente Perfil = "ATENDENTE"
	PerfilMecanico  Perfil = "MECANICO"
)

type Usuario struct {
	ID        int64     `json:"id"`
	Nome      string    `json:"nome" binding:"required"`
	Login     string    `json:"login" binding:"required"`
	SenhaHash string    `json:"-"`
	Perfil    Perfil    `json:"perfil" binding:"required,oneof=GESTOR ATENDENTE MECANICO"`
	Ativo     bool      `json:"ativo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
