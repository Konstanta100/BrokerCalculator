include .env

LOCAL_BIN:=$(CURDIR)/bin

ENV_DIR = .env

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc-gen:
	bin/sqlc generate

migration-status:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v

migration-add:
	goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v


