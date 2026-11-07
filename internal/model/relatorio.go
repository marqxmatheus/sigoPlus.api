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

// FiltroOrdens usa Inicio como limite inclusivo e Fim como limite exclusivo.
type FiltroOrdens struct {
	Inicio    *time.Time
	Fim       *time.Time
	ClienteID *int64
}

type ItemRelatorioOS struct {
	OrdemServicoID int64           `json:"ordem_servico_id"`
	DataAbertura   time.Time       `json:"data_abertura"`
	Situacao       SituacaoOS      `json:"situacao"`
	ClienteID      int64           `json:"cliente_id"`
	ClienteNome    string          `json:"cliente_nome"`
	VeiculoPlaca   string          `json:"veiculo_placa"`
	VeiculoMarca   string          `json:"veiculo_marca"`
	VeiculoModelo  string          `json:"veiculo_modelo"`
	ValorTotal     float64         `json:"valor_total"`
	FormaPagamento *FormaPagamento `json:"forma_pagamento"`
	DataPagamento  *time.Time      `json:"data_pagamento"`
}

type RelatorioOSPeriodo struct {
	DataEmissao time.Time         `json:"data_emissao"`
	DataInicio  string            `json:"data_inicio,omitempty"`
	DataFim     string            `json:"data_fim,omitempty"`
	Itens       []ItemRelatorioOS `json:"itens"`
	Quantidade  int               `json:"quantidade"`
	ValorTotal  float64           `json:"valor_total"`
}

type RelatorioOSCliente struct {
	DataEmissao time.Time         `json:"data_emissao"`
	Cliente     *Cliente          `json:"cliente,omitempty"`
	Itens       []ItemRelatorioOS `json:"itens"`
	Quantidade  int               `json:"quantidade"`
	ValorTotal  float64           `json:"valor_total"`
}
