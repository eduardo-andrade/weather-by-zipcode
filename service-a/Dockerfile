# syntax=docker/dockerfile:1.4

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependência da raiz
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o projeto
COPY . .

# Compila o serviço A (binário a partir de ./service-a)
RUN go build -o /service-a ./service-a

# Imagem final mínima
FROM alpine:latest

# Copia apenas o binário compilado
COPY --from=builder /service-a /service-a

EXPOSE 8080

# Executa o serviço
ENTRYPOINT ["/service-a"]
