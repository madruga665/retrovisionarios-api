# Retrovisionários API

API Backend desenvolvida em **Go** para o projeto Retrovisionários. O sistema gerencia eventos e informações da banda.

## 🚀 Tecnologias

- **Linguagem:** [Go](https://go.dev/) (1.25+)
- **Framework Web:** [Gin Gonic](https://github.com/gin-gonic/gin)
- **Database:** PostgreSQL
- **Driver SQL:** [pgx/v5](https://github.com/jackc/pgx) (High performance)
- **Gerenciamento de Dependências:** Go Modules

## ARCHITECTURE

O projeto evoluiu para uma estrutura baseada em **Domain-Driven Design (DDD)** e **Versionamento**, facilitando a escalabilidade e manutenção:

- **cmd/server**: Ponto de entrada (Main, Wiring de dependências).
- **internal/app/v1**: Contém a lógica da versão 1 da API.
    - **routes.go**: Definição de rotas e grupos da V1.
    - **events/**: Domínio de Eventos.
        - **controllers**: Handler HTTP.
        - **services**: Regras de negócio.
        - **repositories**: Acesso a dados (SQL).
        - **models**: Estruturas de dados.
- **internal/db**: Configuração de infraestrutura (Conexão DB).
- **config/env**: Gerenciamento de variáveis de ambiente.

## 🛠️ Como Rodar Localmente

### Pré-requisitos

- Go instalado
- Docker e Docker Compose (para o banco de dados)

### 1. Configurar Banco de Dados

O projeto possui um `docker-compose.yml` que sobe o PostgreSQL e já executa o script de inicialização (`sql/init.sql`).

```bash
docker-compose up -d
```

### 2. Configurar Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com base no exemplo abaixo (ou use as configurações padrão do docker-compose):

```bash
DATABASE_URL=postgres://postgres:password@localhost:5432/retrovisionarios
PORT=5000
```

### 3. Executar a API

#### Com Hot Reload (Recomendado para Desenvolvimento)

O projeto utiliza o [Air](https://github.com/air-verse/air) para recompilação automática:

```bash
# Instalar o air se ainda não tiver
go install github.com/air-verse/air@latest

# Rodar com hot reload
air
```

#### Execução Padrão

```bash
# Baixar dependências
go mod tidy

# Rodar o servidor
go run cmd/server/main.go
```

A API estará rodando em: `http://localhost:5000`

## 🔌 API Endpoints

### Events

| Método | Rota         | Descrição                          |
| ------ | ------------ | ---------------------------------- |
| `GET`  | `/v1/events` | Lista todos os eventos cadastrados |

#### Exemplo de Resposta (GET /v1/events)

```json
{
  "result": [
    {
      "id": 1,
      "date": "2025-12-08T18:30:00Z",
      "name": "Aniversário Moto Club Dragões",
      "flyer": "http://aws.bucket.com/foto/1"
    }
  ]
}
```

## 🧪 Testes e Ferramentas
