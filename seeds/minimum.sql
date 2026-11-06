BEGIN;

-- Cliente com todos os dados necessários para o cadastro e o atendimento.
INSERT INTO clientes (nome, cpf_cnpj, telefone, email, endereco)
SELECT
    'Cliente padrão SIGO+',
    '52998224725',
    '44999990000',
    'cliente@sigo.local',
    'Rua das Oficinas, 100 - Centro'
WHERE NOT EXISTS (
    SELECT 1
    FROM clientes
    WHERE cpf_cnpj = '52998224725'
      AND ativo = TRUE
);

INSERT INTO clientes (nome, cpf_cnpj, telefone, email, endereco)
SELECT
    'Matheus Caetano Marques',
    '11324207922',
    '44988624646',
    'matheus.marques5215@gmail.com',
    'Rua Sapucaí 412, Cianorte, Paraná'
WHERE NOT EXISTS (
    SELECT 1
    FROM clientes
    WHERE cpf_cnpj = '11324207922'
      AND ativo = TRUE
);

-- Veículo vinculado ao cliente padrão.
INSERT INTO veiculos (cliente_id, placa, marca, modelo, ano)
SELECT
    id,
    'SIG0A26',
    'Honda',
    'CG 160 Titan',
    2025
FROM clientes
WHERE cpf_cnpj = '52998224725'
  AND ativo = TRUE
ORDER BY id DESC
LIMIT 1
ON CONFLICT (placa) DO NOTHING;

INSERT INTO veiculos (cliente_id, placa, marca, modelo, ano)
SELECT
    id,
    'ABC1234',
    'Honda',
    'Biz',
    2010
FROM clientes
WHERE cpf_cnpj = '11324207922'
  AND ativo = TRUE
ORDER BY id DESC
LIMIT 1
ON CONFLICT (placa) DO NOTHING;

-- Serviço que pode ser selecionado na criação da ordem.
INSERT INTO servicos (descricao, valor_padrao, tempo_estimado_min)
VALUES ('Revisão preventiva', 150.00, 90)
ON CONFLICT (descricao) DO NOTHING;

INSERT INTO servicos (descricao, valor_padrao, tempo_estimado_min)
VALUES ('Troca de Óleo', 50.00, 15)
ON CONFLICT (descricao) DO NOTHING;

-- Peça com preço, estoque atual e estoque mínimo preenchidos.
INSERT INTO pecas (
    codigo,
    descricao,
    valor_unitario,
    qtd_estoque,
    estoque_minimo
)
VALUES ('OLEO-10W30', 'Óleo do motor 10W30', 39.90, 20, 5)
ON CONFLICT (codigo) DO NOTHING;

INSERT INTO pecas (
    codigo,
    descricao,
    valor_unitario,
    qtd_estoque,
    estoque_minimo
)
VALUES ('OLEO', 'Óleo Sintético', 20.00, 5, 10)
ON CONFLICT (codigo) DO NOTHING;

COMMIT;
