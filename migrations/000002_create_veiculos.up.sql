CREATE TABLE veiculos (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    cliente_id  BIGINT       NOT NULL REFERENCES clientes (id) ON DELETE RESTRICT,
    placa       VARCHAR(7)   NOT NULL UNIQUE,
    marca       VARCHAR(50)  NOT NULL DEFAULT '',
    modelo      VARCHAR(50)  NOT NULL DEFAULT '',
    ano         INTEGER      NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT ck_veiculos_placa CHECK (placa ~ '^[A-Z]{3}[0-9][A-Z0-9][0-9]{2}$')
);

CREATE INDEX idx_veiculos_cliente_id ON veiculos (cliente_id);