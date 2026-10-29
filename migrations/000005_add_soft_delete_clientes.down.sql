DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM clientes WHERE ativo = FALSE) THEN
        RAISE EXCEPTION 'não é possível reverter: existem clientes com dados anonimizados';
    END IF;
END $$;

DROP INDEX uq_clientes_cpf_cnpj_ativo;

ALTER TABLE clientes
    ADD CONSTRAINT clientes_cpf_cnpj_key UNIQUE (cpf_cnpj),
    DROP COLUMN deleted_at,
    DROP COLUMN ativo;
