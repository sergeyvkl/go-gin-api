# Dockerfile
# Этап 1: Сборка
FROM golang:1.22.2-alpine AS builder

# Установка необходимых инструментов
RUN apk add --no-cache git ca-certificates

# Установка рабочей директории
WORKDIR /app/go-gin-api

# Копирование go.mod и go.sum
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main .

# Этап 2: Запуск
FROM alpine:latest

# Установка CA сертификатов для HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Установка часового пояса (опционально)
ENV TZ=Europe/Moscow

# Создание непривилегированного пользователя
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /app/go-gin-api

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/go-gin-api/main /app/go-gin-api/main

# Копирование файла .env (если есть)
COPY --from=builder /app/go-gin-api/.env /app/go-gin-api/.env

# Переключение на непривилегированного пользователя
USER appuser

# Рабочая директория
WORKDIR /app/go-gin-api

# Экспорт порта
EXPOSE 8080

# Запуск приложения
CMD ["/app/go-gin-api/main"]