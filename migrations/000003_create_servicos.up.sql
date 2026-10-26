CREATE TABLE servicos (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    descricao           VARCHAR(100)  NOT NULL UNIQUE,
    valor_padrao        NUMERIC(10,2) NOT NULL DEFAULT 0 CHECK (valor_padrao >= 0),
    tempo_estimado_min  INTEGER       NOT NULL DEFAULT 0 CHECK (tempo_estimado_min >= 0),
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT now()
);