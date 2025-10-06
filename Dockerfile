FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mainaa ./cmd/short
RUN CGO_ENABLED=0 GOOS=linux go build -o migrationsfa ./cmd/DirMig

FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY --from=builder /app/mainaa .
COPY --from=builder /app/migrationsfa .
RUN chmod +x mainaa migrationsfa

COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/scripts/start-with-migrations.sh .
RUN chmod +x start-with-migrations.sh

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./start-with-migrations.sh"]