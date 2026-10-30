CREATE TABLE ordens_servico (
    id               BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cliente_id       BIGINT        NOT NULL REFERENCES clientes (id) ON DELETE RESTRICT,
    veiculo_id       BIGINT        NOT NULL REFERENCES veiculos (id) ON DELETE RESTRICT,
    situacao         VARCHAR(20)   NOT NULL DEFAULT 'ORCADO',
    valor_total      NUMERIC(10,2) NOT NULL DEFAULT 0,
    observacao       VARCHAR(500)  NOT NULL DEFAULT '',
    forma_pagamento  VARCHAR(10),
    data_pagamento   TIMESTAMPTZ,
    created_at       TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ   NOT NULL DEFAULT now(),

    CONSTRAINT ck_os_situacao CHECK (situacao IN
        ('ORCADO', 'EM_FILA', 'EM_ANDAMENTO', 'AGUARDANDO_PECA', 'FINALIZADO', 'PAGO')),
    CONSTRAINT ck_os_forma_pagamento CHECK (forma_pagamento IN ('DINHEIRO', 'PIX', 'CARTAO'))
);

CREATE INDEX idx_os_cliente_id ON ordens_servico (cliente_id);
CREATE INDEX idx_os_veiculo_id ON ordens_servico (veiculo_id);
CREATE INDEX idx_os_situacao   ON ordens_servico (situacao);

CREATE TABLE os_servicos (
    id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ordem_servico_id  BIGINT        NOT NULL REFERENCES ordens_servico (id) ON DELETE CASCADE,
    servico_id        BIGINT        NOT NULL REFERENCES servicos (id) ON DELETE RESTRICT,
    valor             NUMERIC(10,2) NOT NULL CHECK (valor >= 0)
);

CREATE INDEX idx_os_servicos_os ON os_servicos (ordem_servico_id);

CREATE TABLE os_pecas (
    id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    ordem_servico_id  BIGINT        NOT NULL REFERENCES ordens_servico (id) ON DELETE CASCADE,
    peca_id           BIGINT        NOT NULL REFERENCES pecas (id) ON DELETE RESTRICT,
    quantidade        INTEGER       NOT NULL CHECK (quantidade > 0),
    valor_unitario    NUMERIC(10,2) NOT NULL CHECK (valor_unitario >= 0)
);

CREATE INDEX idx_os_pecas_os ON os_pecas (ordem_servico_id);