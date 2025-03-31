-- name: UserCreate :one
insert into investing.users(name, email)
VALUES ($1, $2) returning id;

-- name: UserById :one
select * from investing.users where id = $1;

-- name: AccountById :one
select * from investing.accounts where id = $1;

-- name: AccountsByUserId :many
select * from investing.accounts where user_id = $1;

-- name: AccountCreate :one
INSERT INTO investing.accounts (
    id,
    user_id,
    name,
    status,
    type,
    access_level,
    opened_date,
    closed_date
) VALUES ($1,$2,$3,$4, $5,$6,$7,$8)
RETURNING id;

-- name: OperationsByInstrument :many
select * from investing.operations where account_id = $1 and instrument_type = $2;

-- name: OperationsByAccountId :many
SELECT * FROM investing.operations WHERE account_id = $1;

-- name: OperationsByInstrumentAndDateRange :many
SELECT * FROM investing.operations
WHERE account_id = $1
  AND instrument_type = $2
  AND date BETWEEN $3 AND $4;

-- name: OperationCreate :one
INSERT INTO investing.operations (
    id,
    figi,
    instrument_type,
    quantity,
    payment,
    currency,
    date,
    account_id
) VALUES ($1, $2, $3, $4, $5,  $6, $7, $8)
    RETURNING id;

-- name: BulkInsertOperations :copyfrom
INSERT INTO investing.operations (
    id,
    figi,
    instrument_type,
    quantity,
    payment,
    currency,
    date,
    account_id
) VALUES ($1,$2,$3,$4, $5,  $6, $7, $8);

