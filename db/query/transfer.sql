-- name: CreateTransfer :one
INSERT INTO transfers (
  from_ccount_id,
  to_ccount_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfer :one
SELECT * from transfers
where id = $1 LIMIT 1;