ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o app main.go

# Указываем порт, который будет использовать приложение
EXPOSE 8182

# Команда для запуска приложения
CMD ["./app"]