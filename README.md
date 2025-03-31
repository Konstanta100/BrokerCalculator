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

### Тестирование
Написаны юнит-тесты на core логику приложения. Плюсом будут тесты на транспортном уровне и на уровне хранения.