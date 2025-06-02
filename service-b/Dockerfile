# syntax=docker/dockerfile:1.4

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copia o go.mod e go.sum da raiz
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte do projeto
COPY . .

# Compila o serviço B (entrypoint em service-b/main.go)
RUN go build -o /service-b ./service-b

# Imagem final
FROM alpine:latest

# Copia o binário compilado
COPY --from=builder /service-b /service-b

# Define o binário como ponto de entrada
ENTRYPOINT ["/service-b"]
