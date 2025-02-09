# Используем официальный образ Golang
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

# Копируем файлы проекта
COPY main.go .

# Собираем бинарник
RUN go build -o pinger

# Устанавливаем Docker CLI
RUN apt-get update && apt-get install -y docker.io

# Запускаем сервис
CMD ["./pinger"]
