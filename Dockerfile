# Используем базовый образ Golang
FROM golang:1.23-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Устанавливаем необходимые инструменты для CGO
RUN apk add --no-cache gcc musl-dev

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Включаем CGO и компилируем проект
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o main .


# Указываем порт
EXPOSE 8080

# Запуск приложения
CMD ["./main"]
