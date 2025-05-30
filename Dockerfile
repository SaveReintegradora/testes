# Etapa de build
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o minha-api

# Imagem final
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/minha-api .
COPY .env .env

EXPOSE 8080

CMD ["./minha-api"]