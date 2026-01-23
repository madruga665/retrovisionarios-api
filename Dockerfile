# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Instalar dependências necessárias para build (se houver CGO, mas aqui vamos desabilitar)
# RUN apk add --no-cache git

# Copiar arquivos de dependência primeiro para cachear camadas
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o binário
# CGO_ENABLED=0 garante um binário estático puro (sem deps de sistema)
# -ldflags="-w -s" remove símbolos de debug para diminuir o tamanho
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Stage 2: Runner
FROM alpine:latest

WORKDIR /app

# Instalar certificados CA para chamadas HTTPS externas e timezone
RUN apk --no-cache add ca-certificates tzdata

# Configurar Timezone (opcional, mas recomendado para BR)
ENV TZ=America/Sao_Paulo

# Copiar apenas o binário da etapa anterior
COPY --from=builder /app/server .

# Expor a porta
EXPOSE 5000

# Variáveis de ambiente padrão para prod
ENV GIN_MODE=release

# Executar
CMD ["./server"]
