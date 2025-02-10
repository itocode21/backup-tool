FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o backup-tool ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /app/backup-tool .

# Копирование конфигурационных файлов
COPY pkg/config/test_config_mysql.yaml /app/config/mysql.yaml
COPY pkg/config/test_config_postgresql.yaml /app/config/postgresql.yaml
COPY pkg/config/test_config_mongodb.yaml /app/config/mongodb.yaml

RUN mkdir -p /app/data/logs && mkdir -p /app/data/backups

ENTRYPOINT ["/app/backup-tool"]