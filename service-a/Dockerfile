FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /service-a ./service-a

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /service-a /service-a

EXPOSE 8080

ENTRYPOINT ["/service-a"]
