// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: transfer.sql

package db

import (
	"context"
)

const createTransfer = `-- name: CreateTransfer :one
INSERT INTO transfers (
  from_ccount_id,
  to_ccount_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING id, from_ccount_id, to_ccount_id, amount, created_at
`

type CreateTransferParams struct {
	FromCcountID int64 `json:"from_ccount_id"`
	ToCcountID   int64 `json:"to_ccount_id"`
	Amount       int64 `json:"amount"`
}

func (q *Queries) CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, createTransfer, arg.FromCcountID, arg.ToCcountID, arg.Amount)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromCcountID,
		&i.ToCcountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getTransfer = `-- name: GetTransfer :one
SELECT id, from_ccount_id, to_ccount_id, amount, created_at from transfers
where id = $1 LIMIT 1
`

func (q *Queries) GetTransfer(ctx context.Context, id int64) (Transfer, error) {
	row := q.db.QueryRowContext(ctx, getTransfer, id)
	var i Transfer
	err := row.Scan(
		&i.ID,
		&i.FromCcountID,
		&i.ToCcountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
