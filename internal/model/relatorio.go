package model

import "time"

type FiltroEstoque struct {
	Busca              string
	ApenasAbaixoMinimo bool
}

type ItemRelatorioEstoque struct {
	PecaID        int64   `json:"peca_id"`
	Codigo        string  `json:"codigo"`
	Descricao     string  `json:"descricao"`
	QtdEstoque    int     `json:"qtd_estoque"`
	EstoqueMinimo int     `json:"estoque_minimo"`
	ValorUnitario float64 `json:"valor_unitario"`
	ValorTotal    float64 `json:"valor_total"`
	AbaixoMinimo  bool    `json:"abaixo_minimo"`
}

type RelatorioEstoque struct {
	DataEmissao       time.Time              `json:"data_emissao"`
	Itens             []ItemRelatorioEstoque `json:"itens"`
	TotalPecas        int                    `json:"total_pecas"`
	QuantidadeTotal   int                    `json:"quantidade_total"`
	ValorTotal        float64                `json:"valor_total"`
	TotalAbaixoMinimo int                    `json:"total_abaixo_minimo"`
}