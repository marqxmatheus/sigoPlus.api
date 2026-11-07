# SIGO+ API

## Ambiente completo com Docker

O passo a passo completo para Windows e Linux está em
[`GUIA_DOCKER.md`](GUIA_DOCKER.md).

Esta é a forma recomendada de executar o projeto em uma máquina nova. É
necessário ter apenas Docker Engine com Compose no Linux ou Docker Desktop com
WSL2 no Windows.

Os repositórios devem estar lado a lado e manter estes nomes:

```text
Sigo/
├── sigoPlus.api/
└── sigoPlus.web/
```

No diretório `sigoPlus.api`, execute:

```bash
docker compose up --build
```

O Compose cria o banco `sigo`, aplica todas as migrations, cadastra os dados
mínimos para criar uma ordem de serviço e inicia:

- Web: http://localhost:5173
- API: http://localhost:8080
- PostgreSQL: localhost:5432

Os arquivos locais são montados nos containers. O Vite atualiza o navegador
automaticamente quando o web é alterado. Após alterações na API, reinicie
somente esse serviço:

```bash
docker compose restart api
```

Para acompanhar os logs:

```bash
docker compose logs -f api web
```

Para encerrar os serviços sem apagar o banco:

```bash
docker compose down
```

Para também apagar o volume e todos os dados do banco, use com cuidado:

```bash
docker compose down -v
```

As portas e credenciais locais podem ser personalizadas copiando
`.env.example` para `.env`. Em caso de conflito com um PostgreSQL já instalado,
por exemplo, altere `POSTGRES_PORT` para `5433`.

## Banco de dados local

Requisitos: PostgreSQL, `psql`, Docker e Make.

A configuração padrão usa PostgreSQL em `localhost:5432`, usuário e senha
`postgres` e banco `sigo`.

```bash
make db-create
make migrate-up
make seed-minimum
make migrate-version
```

O seed mínimo é idempotente: ele cadastra duas combinações completas de
cliente, veículo, serviço e peça apenas quando esses registros ainda não
existem. Por isso, o comando pode ser executado novamente sem gerar
duplicidades.

As migrations ficam em `migrations/` e são executadas pelo
[`golang-migrate`](https://github.com/golang-migrate/migrate) em um container.
O banco registra a versão aplicada na tabela `schema_migrations`.

Para desfazer somente a migration mais recente:

```bash
make migrate-down
```

As URLs podem ser substituídas na chamada do Make:

```bash
make migrate-up DATABASE_URL='postgres://usuario:senha@localhost:5432/sigo?sslmode=disable'
```

## Executar a API

Com Go instalado:

```bash
go run ./cmd/api
```

As variáveis aceitas estão documentadas em `.env.example`. A aplicação lê as
variáveis do ambiente; o arquivo `.env` não é carregado automaticamente.

## Relatórios

| Relatório | Endpoint | Filtros |
|---|---|---|
| Estoque de peças (RF7) | `/api/v1/relatorios/estoque` | `busca`, `abaixo_minimo` |
| Ordens de serviço por período | `/api/v1/relatorios/ordens-periodo` | `data_inicio`, `data_fim` |
| Ordens de serviço por cliente | `/api/v1/relatorios/ordens-cliente` | `cliente_id` |

### Estoque de peças

- `busca`: procura no código ou na descrição, ignorando maiúsculas e minúsculas (não ignora acentos).
- `abaixo_minimo=true`: só as peças com `qtd_estoque < estoque_minimo`.
- Cada item tem `abaixo_minimo: true/false`, para destacar a linha na tela.

### Ordens por período

- Usa a **data de abertura** da ordem e traz **qualquer situação**.
- `data_inicio` e `data_fim` são inclusivos: `2026-11-07` até `2026-11-07` traz o dia 7 inteiro.
- Sem filtros, traz todas as ordens.

### Ordens por cliente

- Com `cliente_id`, a resposta traz o objeto `cliente` (mesmo formato de `GET /clientes/:id`) para o cabeçalho do relatório.
- Sem `cliente_id`, traz todas as ordens e o campo `cliente` não vem.
- Cliente existente sem ordens: `200` com `itens: []`.
Os dois relatórios de ordens vêm ordenados da mais recente para a mais antiga.

### Erros

Mesmo formato do resto da API: `{ "erro": "mensagem" }`.
- `400`: data fora do formato, data inicial maior que a final, `cliente_id` ou `abaixo_minimo` inválidos.
- `404`: `cliente_id` de cliente inexistente.

### Datas

- `data_emissao` vem no horário de Brasília (`-03:00`); `data_abertura` e `data_pagamento` vêm em UTC (`Z`).
- Para exibir, converter sempre com `new Date(valor).toLocaleString('pt-BR')`, que ajusta o fuso automaticamente.

### Tipos sugeridos (`src/types/relatorio.ts`)

```ts
import type { Cliente } from './cliente'
import type { FormaPagamento, SituacaoOS } from './ordemServico'
export type ItemRelatorioEstoque = {
  peca_id: number
  codigo: string
  descricao: string
  qtd_estoque: number
  estoque_minimo: number
  valor_unitario: number
  valor_total: number
  abaixo_minimo: boolean
}
export type RelatorioEstoque = {
  data_emissao: string
  itens: ItemRelatorioEstoque[]
  total_pecas: number
  quantidade_total: number
  valor_total: number
  total_abaixo_minimo: number
}
export type ItemRelatorioOS = {
  ordem_servico_id: number
  data_abertura: string
  situacao: SituacaoOS
  cliente_id: number
  cliente_nome: string
  veiculo_placa: string
  veiculo_marca: string
  veiculo_modelo: string
  valor_total: number
  forma_pagamento: FormaPagamento | null
  data_pagamento: string | null
}
export type RelatorioOSPeriodo = {
  data_emissao: string
  data_inicio?: string
  data_fim?: string
  itens: ItemRelatorioOS[]
  quantidade: number
  valor_total: number
}
export type RelatorioOSCliente = {
  data_emissao: string
  cliente?: Cliente
  itens: ItemRelatorioOS[]
  quantidade: number
  valor_total: number
}
```

### Funções sugeridas (`src/api/relatorios.ts`)

O axios monta a query string a partir de `params` e ignora os valores `undefined`, então os filtros vazios simplesmente não são enviados.
```ts
import type {
  RelatorioEstoque,
  RelatorioOSCliente,
  RelatorioOSPeriodo,
} from '../types/relatorio'
import { api } from './client'
const caminho = '/api/v1/relatorios'
export async function relatorioEstoque(filtros: {
  busca?: string
  abaixo_minimo?: boolean
} = {}) {
  const { data } = await api.get<RelatorioEstoque>(`${caminho}/estoque`, {
    params: filtros,
  })
  return data
}
export async function relatorioOrdensPeriodo(filtros: {
  data_inicio?: string
  data_fim?: string
} = {}) {
  const { data } = await api.get<RelatorioOSPeriodo>(
    `${caminho}/ordens-periodo`,
    { params: filtros },
  )
  return data
}
export async function relatorioOrdensCliente(cliente_id?: number) {
  const { data } = await api.get<RelatorioOSCliente>(
    `${caminho}/ordens-cliente`,
    { params: { cliente_id } },
  )
  return data
}
```

### Exemplo de resposta (`/ordens-cliente?cliente_id=2`)

```json
{
  "data_emissao": "2026-11-07T11:55:32-03:00",
  "cliente": {
    "id": 2,
    "nome": "Matheus Caetano Marques",
    "cpf_cnpj": "11324207922",
    "telefone": "44988624646",
    "email": "matheus.marques5215@gmail.com",
    "endereco": "Rua Sapucaí 412, Cianorte, Paraná",
    "created_at": "...",
    "updated_at": "..."
  },
  "itens": [
    {
      "ordem_servico_id": 1,
      "data_abertura": "2026-11-07T10:34:19Z",
      "situacao": "ORCADO",
      "cliente_id": 2,
      "cliente_nome": "Matheus Caetano Marques",
      "veiculo_placa": "ABC1234",
      "veiculo_marca": "Honda",
      "veiculo_modelo": "Biz",
      "valor_total": 239.9,
      "forma_pagamento": null,
      "data_pagamento": null
    }
  ],
  "quantidade": 1,
  "valor_total": 239.9
}
``` 