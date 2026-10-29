# Guia do ambiente Docker

Este guia explica como executar o SIGO+ em uma máquina nova usando Docker. O
ambiente inclui PostgreSQL, migrations, API e web, sem exigir instalações locais
de Go, Node.js ou PostgreSQL.

## 1. Pré-requisitos

### Windows

1. Instale o [Docker Desktop](https://www.docker.com/products/docker-desktop/).
2. Durante a instalação, mantenha a opção de usar WSL2 habilitada.
3. Abra o Docker Desktop e aguarde até ele indicar que o Docker está em
   execução.

### Linux

Instale Docker Engine e o plugin Docker Compose pelos pacotes recomendados para
sua distribuição. Certifique-se de que o serviço do Docker esteja em execução.

Em qualquer sistema, confirme a instalação em um terminal:

```bash
docker --version
docker compose version
```

## 2. Organização dos repositórios

Os dois repositórios devem ficar lado a lado e conservar estes nomes:

```text
Sigo/
├── sigoPlus.api/
└── sigoPlus.web/
```

Exemplo de clonagem:

```bash
mkdir Sigo
cd Sigo
git clone <URL_DO_REPOSITORIO_DA_API> sigoPlus.api
git clone <URL_DO_REPOSITORIO_DO_WEB> sigoPlus.web
```

O arquivo `compose.yaml` fica no repositório da API e usa o diretório irmão
`../sigoPlus.web` para construir o frontend.

## 3. Primeira execução

Entre no diretório da API.

No Windows PowerShell:

```powershell
cd C:\caminho\para\Sigo\sigoPlus.api
```

No Linux:

```bash
cd ~/caminho/para/Sigo/sigoPlus.api
```

Construa as imagens e inicie o ambiente:

```bash
docker compose up --build
```

Na primeira execução, o Docker irá:

1. Baixar as imagens de PostgreSQL, Go e Node.js.
2. Criar o banco de dados `sigo`.
3. Executar todas as migrations disponíveis.
4. Compilar e iniciar a API.
5. Instalar as dependências e iniciar o web.

O primeiro build pode levar alguns minutos. Aguarde as mensagens indicando que
a API está ouvindo na porta `8080` e que o Vite está pronto.

## 4. Endereços locais

Com os serviços ativos:

- Web: http://localhost:5173
- API: http://localhost:8080
- Healthcheck da API: http://localhost:8080/api/v1/health
- PostgreSQL: `localhost:5432`

## 5. Execução em segundo plano

Para liberar o terminal, use:

```bash
docker compose up --build -d
```

Consulte o estado dos serviços:

```bash
docker compose ps
```

Veja os logs da API e do web:

```bash
docker compose logs -f api web
```

Para sair da visualização dos logs sem encerrar os containers, pressione
`Ctrl+C`.

## 6. Uso durante o desenvolvimento

Os diretórios dos dois projetos são montados nos containers.

- Alterações no web são detectadas automaticamente pelo Vite.
- Depois de alterar código Go, reinicie a API:

```bash
docker compose restart api
```

- Depois de alterar dependências, Dockerfiles ou o próprio Compose, reconstrua:

```bash
docker compose up --build
```

As migrations são executadas antes da API. Ao adicionar uma nova migration,
recrie o fluxo dos serviços com:

```bash
docker compose down
docker compose up --build
```

## 7. Encerrar o ambiente

Quando o Compose estiver rodando em primeiro plano, pressione `Ctrl+C`. Para
remover os containers e a rede, sem apagar o banco:

```bash
docker compose down
```

Os dados permanecem no volume `sigo-plus_postgres_data` e estarão disponíveis
na próxima execução.

Para remover containers e também apagar definitivamente todos os dados do
banco:

```bash
docker compose down -v
```

> Atenção: `docker compose down -v` é destrutivo e não pode ser desfeito sem um
> backup.

## 8. Personalizar portas e credenciais

Os valores padrão funcionam sem arquivo adicional. Para personalizá-los, copie
o exemplo de ambiente na raiz da API.

No Windows PowerShell:

```powershell
Copy-Item .env.example .env
```

No Linux:

```bash
cp .env.example .env
```

As principais opções são:

```dotenv
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=sigo
POSTGRES_PORT=5432
API_HTTP_PORT=8080
WEB_PORT=5173
```

Se a porta `5432` já estiver sendo usada por um PostgreSQL local, por exemplo,
altere somente a porta publicada:

```dotenv
POSTGRES_PORT=5433
```

As conexões entre containers continuam usando as portas internas corretas.

## 9. Acessar o PostgreSQL

Abra o terminal do PostgreSQL dentro do container:

```bash
docker compose exec postgres psql -U postgres -d sigo
```

Para sair do `psql`:

```text
\q
```

Verifique a versão das migrations:

```bash
docker compose exec postgres psql -U postgres -d sigo -c "SELECT * FROM schema_migrations;"
```

## 10. Solução de problemas

### Uma porta já está em uso

Altere `POSTGRES_PORT`, `API_HTTP_PORT` ou `WEB_PORT` no arquivo `.env` e tente
novamente.

### Um serviço não inicia

Confira o estado e os logs:

```bash
docker compose ps -a
docker compose logs api
docker compose logs postgres
docker compose logs migrations
docker compose logs web
```

### Dependências ou imagens desatualizadas

Reconstrua sem usar o cache:

```bash
docker compose build --no-cache
docker compose up
```

### Banco local descartável para recomeçar

Somente se os dados puderem ser apagados:

```bash
docker compose down -v
docker compose up --build
```
