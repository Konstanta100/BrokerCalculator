-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS investing;

CREATE TABLE investing.users
(
    id       uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name     TEXT         NOT NULL,
    email    VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE investing.accounts (
    id VARCHAR(255) PRIMARY KEY,
    user_id uuid NOT NULL,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    access_level VARCHAR(255) NOT NULL,
    opened_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_date TIMESTAMP NULL,
    FOREIGN KEY (user_id) REFERENCES investing.users(id) ON DELETE CASCADE
);

CREATE TABLE investing.operations (
    id VARCHAR(255) PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL,
    figi VARCHAR(255) NOT NULL,
    instrument_type VARCHAR(255) NOT NULL,
    payment DECIMAL(15, 2) NOT NULL,
    quantity BIGINT NOT NULL,
    currency VARCHAR(20) NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES investing.accounts(id) ON DELETE CASCADE
);

CREATE INDEX idx_accounts_user_id ON investing.accounts(user_id);
CREATE INDEX idx_operations_account_id ON investing.operations(account_id);
CREATE INDEX idx_operations_instrument_type ON investing.operations(instrument_type);
CREATE INDEX idx_operations_account_date ON investing.operations(account_id, date);
CREATE INDEX idx_operations_account_instrument_date ON investing.operations(account_id, instrument_type, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE investing.operations;
DROP TABLE investing.accounts;
DROP TABLE investing.users;

DROP SCHEMA investing;
-- +goose StatementEnd
