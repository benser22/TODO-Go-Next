# syntax=docker/dockerfile:1

# ---------- build: compila el binario con deps cacheadas ----------
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -o /out/todo-api ./cmd/todo-api

# ---------- goose: herramienta de migraciones como binario ----------
FROM golang:1.27-alpine AS goose
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.28.0

# ---------- runtime: imagen mínima, usuario root no necesario ----------
FROM alpine:3.22
RUN apk add --no-cache ca-certificates && adduser -u 10001 -D app
COPY --from=build /out/todo-api /usr/local/bin/todo-api
COPY --from=goose /go/bin/goose /usr/local/bin/goose
COPY db/migrations /migrations
USER app
EXPOSE 8080
# Aplica migraciones y luego cede la señal al binario (exec para SIGTERM limpio).
ENTRYPOINT ["sh", "-c", "goose -dir /migrations postgres \"$DATABASE_URL\" up && exec todo-api"]
