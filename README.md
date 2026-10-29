# SIGO+ API

## Banco de dados local

Requisitos: PostgreSQL, `psql`, Docker e Make.

A configuração padrão usa PostgreSQL em `localhost:5432`, usuário e senha
`postgres` e banco `sigo`.

```bash
make db-create
make migrate-up
make migrate-version
```

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
