ALTER TABLE clientes
    ADD COLUMN ativo BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE clientes DROP CONSTRAINT clientes_cpf_cnpj_key;

CREATE UNIQUE INDEX uq_clientes_cpf_cnpj_ativo
    ON clientes (cpf_cnpj)
    WHERE ativo = TRUE;
