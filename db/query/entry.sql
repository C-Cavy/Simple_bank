-- name: CreateEntry :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntry :one
SELECT * from entries
where id = $1 LIMIT 1;


