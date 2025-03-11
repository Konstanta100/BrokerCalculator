ARG GO_VERSION=1.23

FROM golang:${GO_VERSION}

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o app ./cmd/app

# Указываем порт, который будет использовать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./app"]