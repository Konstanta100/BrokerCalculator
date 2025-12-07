ARG GO_VERSION=1.23.4

FROM golang:${GO_VERSION}-alpine as builder

# Установите git и другие необходимые инструменты
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /myapp

FROM alpine:latest

COPY --from=builder /myapp /myapp
COPY .env .env

EXPOSE 8182

CMD ["/myapp"]