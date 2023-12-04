# Указываем базовый образ, содержащий Go
FROM golang:latest

# Копируем файлы проекта в контейнер
COPY . /app

# Устанавливаем рабочую директорию
WORKDIR /app

# Собираем проект
RUN go build -o .\cmd\url-shortener\main.go .

# Определяем команду запуска сервера при запуске контейнера
CMD ["./main"]

