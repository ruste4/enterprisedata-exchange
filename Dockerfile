# Многоэтапная сборка
# Этап 1: Сборка приложения
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git ca-certificates tzdata

# Создание пользователя для приложения
RUN adduser -D -g '' appuser

# Установка рабочей директории
WORKDIR /build

# Копирование go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download
RUN go mod verify

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app cmd/server/main.go

# Этап 2: Финальный образ
FROM scratch

# Импорт временной зоны из builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Импорт CA сертификатов из builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Импорт пользователя и группы из builder
COPY --from=builder /etc/passwd /etc/passwd

# Копирование собранного приложения
COPY --from=builder /build/app /app

# Копирование конфигурационных файлов
COPY --from=builder /build/configs /configs

# Использование непривилегированного пользователя
USER appuser

# Открытие порта
EXPOSE 8080

# Запуск приложения
ENTRYPOINT ["/app"] 