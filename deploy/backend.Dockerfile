FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/api /app/api
COPY backend/migrations /app/migrations

EXPOSE 8080

# Команда указана в docker-compose.yml
