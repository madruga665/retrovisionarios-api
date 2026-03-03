# Retrovisionários API

API Backend desenvolvida em **Go** para o projeto Retrovisionários. O sistema gerencia eventos e informações da banda, com foco em performance, segurança e manutenibilidade.

## 🚀 Tecnologias

- **Linguagem:** [Go](https://go.dev/) (1.25+)
- **Framework Web:** [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL
- **Driver SQL:** [pgx/v5](https://github.com/jackc/pgx) (Pool de conexões de alta performance)
- **Documentação:** [Swagger/OpenAPI](https://github.com/swaggo/swag)
- **Logging:** Structured Logging com `log/slog` (Nativo Go 1.21+)
- **Containerização:** Docker e Docker Compose (Multi-stage builds)

## 🏗️ Arquitetura

O projeto segue uma arquitetura em camadas bem definida, facilitando o desacoplamento e a testabilidade:

- **`cmd/server`**: Ponto de entrada da aplicação, configuração do servidor e injeção de dependências.
- **`config/env`**: Centralização e carregamento de variáveis de ambiente.
- **`internal/app/v1`**: Versão 1 da API com estrutura baseada em domínios.
  - **`events/`**: Módulo de gerenciamento de eventos.
    - **`controllers`**: Camada de transporte (HTTP Handlers).
    - **`services`**: Camada de regras de negócio e interfaces.
    - **`repositories`**: Camada de persistência (SQL).
    - **`models`**: Entidades e tipos customizados (ex: `DateTime`).
- **`internal/db`**: Infraestrutura de conexão ao banco de dados.
- **`docs`**: Arquivos gerados automaticamente para a documentação Swagger.

## 🛠️ Como Rodar Localmente

### Pré-requisitos

- Go 1.21+ instalado
- Docker e Docker Compose

### 1. Configurar Banco de Dados

Suba o contêiner do PostgreSQL (já inclui script de inicialização `sql/init.sql`):

```bash
docker-compose up -d
```

### 2. Variáveis de Ambiente

Crie um arquivo `.env` na raiz:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/retrovisionarios
PORT=5000
TRUSTED_PROXIES=127.0.0.1
ALLOWED_ORIGINS=http://localhost:3000,https://retrovisionarios.vercel.app
LOG_LEVEL=debug
```

### 3. Executar a API

#### Com Hot Reload (Recomendado para Desenvolvimento)

O projeto suporta o [Air](https://github.com/air-verse/air) para recompilação automática ao salvar arquivos:

```bash
# Instalar o air se ainda não tiver
go install github.com/air-verse/air@latest

# Rodar com hot reload (utiliza o arquivo .air.toml se presente ou configurações padrão)
air
```

#### Execução Padrão

```bash
# Baixar dependências
go mod tidy

# Rodar o servidor
go run cmd/server/main.go
```

## 📖 Documentação API (Swagger)

A documentação interativa fica disponível em:
`http://localhost:5000/swagger/index.html`

Para atualizar a documentação após mudanças no código:

```bash
go run github.com/swaggo/swag/cmd/swag init -g cmd/server/main.go
```

## 🔌 Principais Endpoints

| Método | Rota         | Descrição                          | Filtros |
| ------ | ------------ | ---------------------------------- | ------- |
| `GET`  | `/v1/events` | Lista todos os eventos cadastrados | `year`  |
| `POST` | `/v1/events` | Cria um novo evento                | N/A     |

## 🧪 Qualidade de Código e Melhores Práticas

- **Context Propagation:** Todas as camadas respeitam o `context.Context` para cancelamento de requisições e timeouts.
- **Graceful Shutdown:** O servidor encerra conexões de forma segura ao receber sinais do SO.
- **CORS:** Proteção configurada via middleware para origens específicas.
- **Type Safety:** Uso de tipos customizados como `DateTime` para garantir formatos consistentes de data/hora no JSON.
- **Testes Unitários:** Cobertura de lógica nos controllers com mocks.

```bash
# Rodar testes
go test ./...
```

## 📦 Docker Deployment

O projeto utiliza **Multi-stage Build** para gerar imagens mínimas (Alpine Linux):

```bash
docker build -t retrovisionarios-api .
docker run -p 5000:5000 retrovisionarios-api
```
