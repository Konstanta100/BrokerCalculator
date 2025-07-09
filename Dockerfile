ARG GO_VERSION=1.23.4

FROM golang:${GO_VERSION}-alpine as builder

# Устанавливаем git и настраиваем DNS через Docker-флаги
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

# Используем альтернативный GOPROXY и отключаем IPv6
ENV GOPROXY=https://goproxy.cn,direct
ENV GOINSECURE=proxy.golang.org
RUN echo "Building with GOPROXY=${GOPROXY}" && \
    go mod download

COPY . .

RUN go build -o /myapp

FROM alpine:latest

COPY --from=builder /myapp /myapp
COPY .env .env

EXPOSE 8182

CMD ["/myapp"]