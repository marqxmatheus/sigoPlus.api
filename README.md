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
