FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Migration tool.
RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Swagger.
COPY "cmd/todo/docs/swagger.yaml" "/app/cmd/todo/docs/swagger.yaml"

# App.
RUN go build -o /app/lazy-todo ./cmd/todo

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/lazy-todo .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder "/app/cmd/todo/docs/swagger.yaml" "/app/cmd/todo/docs/swagger.yaml"

COPY migrations /app/migrations

# Postgres client for migration.
RUN apk add --no-cache curl postgresql-client

ENV DATABASE_URL=postgres://postgres:password@postgres:5432/dev

EXPOSE 8080

# Wait for postgres to be ready and run migration, after that run app.
CMD ["sh", "-c", "until pg_isready -h postgres -p 5432; do echo 'Waiting for PostgreSQL...'; sleep 2; done && migrate -path=/app/migrations -database $DATABASE_URL up && /app/lazy-todo"]