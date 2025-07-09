include .env

LOCAL_BIN:=$(CURDIR)/bin

ENV_DIR = .env

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
	GOBIN=$(LOCAL_BIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0

sqlc-gen:
	bin/sqlc generate

migration-status:
	bin/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} status -v

migration-add:
	bin/goose -dir ${MIGRATION_DIR} create $(name) sql

migration-up:
	bin/goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

migration-down:
	goose -dir ${MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v


