CREATE TABLE pecas (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    codigo          VARCHAR(30)   NOT NULL UNIQUE,
    descricao       VARCHAR(100)  NOT NULL,
    valor_unitario  NUMERIC(10,2) NOT NULL DEFAULT 0 CHECK (valor_unitario >= 0),
    qtd_estoque     INTEGER       NOT NULL DEFAULT 0 CHECK (qtd_estoque >= 0),
    estoque_minimo  INTEGER       NOT NULL DEFAULT 0 CHECK (estoque_minimo >= 0),
    created_at      TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ   NOT NULL DEFAULT now()
);