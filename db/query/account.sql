-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * from accounts
where id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * from accounts
where id = $1 LIMIT 1
FOR no key UPDATE;

-- name: ListAccounts :many
SELECT * from accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
where id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts
where id = $1;