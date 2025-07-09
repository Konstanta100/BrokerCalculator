## BrokerCalculator HTTP-сервер статистики по брокерскому cчету

[![Go Report Card](https://goreportcard.com/badge/github.com/Konstanta100/BrokerCalculator)](https://goreportcard.com/report/github.com/Konstanta100/BrokerCalculator)
[![CI Status](https://github.com/Konstanta100/BrokerCalculator/actions/workflows/ci.yml/badge.svg)](https://github.com/Konstanta100/BrokerCalculator/actions)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.14-blue)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

## Описание:
В проекте представлен HTTP-сервер, который обрабатывает запросы
для подсчёта комиссий у брокера, за различные периоды и по Брокерским счетам пользователя,
в частности Т-Инвестиции

Документация по SDC подключения к брокеру в T-Инвестиции:
https://russianinvestments.github.io/investAPI/api_protocols/

### Требования:
- HTTP-сервер должен хранить информацию по счетам пользователя в СУБД: аккаунты, операции
- Реализовано подключение по средствам gRPC, к серверу Брокера
- Реализованы API методы для выгрузки данных со счетов пользователя, а также логика со сбором статистики в частности подсчёт комиссий за указанный интервал

### Развертывание
Развертывание сервиса должно осуществляться с использованием docker compose в директории с проектом.

## Порядок выполнения API запросов для подсчёта комиссий за период

### Создайте пользователя
POST /api/user/create"
{
  "name": "Test Testov",
  "email": "test@gmail.com"
}

### Скачайте данные по аккаунту
GET /accounts/load 

### Скачайте операции по аккаунту
POST /operations/load

### Подсчёт комиссии за период
POST /operations/commission

## Тестирование
Написаны юнит-тесты на core логику приложения. Плюсом будут тесты на транспортном уровне и на уровне хранения.

## Разворачивание не сервере
1) sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
2) curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
3) echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
4) sudo apt update
5) sudo apt install -y docker-ce docker-ce-cli containerd.io
6) sudo docker --version
7) sudo curl -L "https://github.com/docker/compose/releases/download/v2.23.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
8) sudo chmod +x /usr/local/bin/docker-compose
9) git clone https://github.com/Konstanta100/BrokerCalculator.git
10) cd BrokerCalculator/
11) sudo docker compose down && sudo docker compose up --build -d
12) sudo mkdir -p /etc/docker
    42  echo '{
    "dns": ["8.8.8.8", "1.1.1.1"],
    "ipv6": false
    }' | sudo tee /etc/docker/daemon.json
13) sudo systemctl restart docker
14) wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
15) sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
16) echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
17) source ~/.bashrc
18) sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
19) sudo -E go version
20) sudo make sqlc-gen
21) sudo make migration-status
22) sudo make migration-up
23) sudo make migration-status
24) sudo docker compose down && sudo docker compose up --build -d



