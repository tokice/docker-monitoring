# Базовый образ
FROM golang:1.23

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы
COPY go.mod ./
COPY main.go ./

# Загружаем зависимости
RUN go mod tidy

# Собираем приложение
RUN go build -o app

# Запускаем приложение
CMD ["/app/app"]
