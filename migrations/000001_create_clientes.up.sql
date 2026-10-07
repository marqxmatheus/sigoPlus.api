CREATE TABLE clientes (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nome        VARCHAR(100) NOT NULL,
    cpf_cnpj    VARCHAR(14)  NOT NULL UNIQUE,
    telefone    VARCHAR(20)  NOT NULL DEFAULT '',
    email       VARCHAR(100) NOT NULL DEFAULT '',
    endereco    VARCHAR(200) NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT ck_clientes_cpf_cnpj CHECK (cpf_cnpj ~ '^([0-9]{11}|[0-9]{14})$')
);