# Билд стадия
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем только файлы модулей сначала для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Билдим приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /1337b04rd ./cmd/main.go

# Финальная стадия
FROM alpine:latest

WORKDIR /app

# Копируем бинарник
COPY --from=builder /1337b04rd .
# Копируем миграции
COPY ./migrations ./migrations

EXPOSE 8080

CMD ["./1337b04rd"]