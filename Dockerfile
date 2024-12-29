# этап сборки
FROM golang:1.23.4 AS builder

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

# ЭТАП РАБОЧЕЙ СРЕДЫ (runtime)
FROM ubuntu:20.04

WORKDIR /app

# копируем собранное приложение
COPY --from=builder /my_app .

# копируем директорию веб-ресурсов
COPY web ./web

# создаём директорию для базы данных
RUN mkdir -p /app/database

# открываем порт для приложения
EXPOSE 7540

# запускаем приложение
CMD ["/app/my_app"]