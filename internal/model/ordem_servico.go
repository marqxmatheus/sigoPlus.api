package model

import "time"

type SituacaoOS string

const (
	SituacaoOrcado         SituacaoOS = "ORCADO"
	SituacaoEmFila         SituacaoOS = "EM_FILA"
	SituacaoEmAndamento    SituacaoOS = "EM_ANDAMENTO"
	SituacaoAguardandoPeca SituacaoOS = "AGUARDANDO_PECA"
	SituacaoFinalizado     SituacaoOS = "FINALIZADO"
	SituacaoPago           SituacaoOS = "PAGO"
)

type FormaPagamento string

const (
	PagamentoDinheiro FormaPagamento = "DINHEIRO"
	PagamentoPix      FormaPagamento = "PIX"
	PagamentoCartao   FormaPagamento = "CARTAO"
)

type OrdemServico struct {
	ID             int64                     `json:"id"`
	ClienteID      int64                     `json:"cliente_id" binding:"required"`
	ClienteNome    string                    `json:"cliente_nome"`
	VeiculoID      int64                     `json:"veiculo_id" binding:"required"`
	VeiculoPlaca   string                    `json:"veiculo_placa"`
	Situacao       SituacaoOS                `json:"situacao"`
	ValorTotal     float64                   `json:"valor_total"`
	Observacao     string                    `json:"observacao"`
	Servicos       []OrdemServicoItemServico `json:"servicos,omitempty" binding:"dive"`
	Pecas          []OrdemServicoItemPeca    `json:"pecas,omitempty" binding:"dive"`
	FormaPagamento *FormaPagamento           `json:"forma_pagamento,omitempty"`
	DataPagamento  *time.Time                `json:"data_pagamento,omitempty"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

type OrdemServicoItemServico struct {
	ID             int64   `json:"id"`
	OrdemServicoID int64   `json:"ordem_servico_id"`
	ServicoID      int64   `json:"servico_id" binding:"required"`
	Descricao      string  `json:"descricao"`
	Valor          float64 `json:"valor"`
}

type OrdemServicoItemPeca struct {
	ID             int64   `json:"id"`
	OrdemServicoID int64   `json:"ordem_servico_id"`
	PecaID         int64   `json:"peca_id" binding:"required"`
	Descricao      string  `json:"descricao"`
	Quantidade     int     `json:"quantidade" binding:"required,gt=0"`
	ValorUnitario  float64 `json:"valor_unitario"`
}
