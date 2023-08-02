FROM golang:1.18-alpine
LABEL authors="zatrasz"

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /account

# Копируем файлы проекта внутрь контейнера
COPY go.mod go.sum ./

COPY cmd/main.go ./

COPY ./ ./

# Устанавливаем зависимости
RUN go mod download

# Собираем приложение
RUN go build -o account .

CMD ["./account"]