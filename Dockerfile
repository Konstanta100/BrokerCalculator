ARG GO_VERSION=1.23.4

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app

# Решение проблемы DNS и IPv6
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf && \
    echo "nameserver 1.1.1.1" >> /etc/resolv.conf && \
    apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /myapp

FROM alpine:latest

COPY --from=builder /myapp /myapp
COPY .env .env

EXPOSE 8182

CMD ["/myapp"]