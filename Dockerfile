# этап сборки
FROM golang:1.23.4-alpine AS builder

# устанавливаем необходимые инструменты
RUN apk add --no-cache git

# собираем для ARM64
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64

WORKDIR /app

# копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# копируем исходный код
COPY . .

# собираем Go-приложение
RUN go build -o /my_app ./app/cmd

# этап рабочей среды (runtime)
FROM alpine:latest

WORKDIR /app

# копируем собранное приложение
COPY --from=builder /my_app .

# копируем директорию веб-ресурсов
COPY web ./web

# создаём директорию для базы данных
RUN mkdir -p /app/database

# копируем .env файл
COPY .env .env

# запускаем приложение
CMD ["sh", "-c", "/app/my_app"]