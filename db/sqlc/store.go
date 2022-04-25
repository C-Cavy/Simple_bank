package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// Store provides all functions to execute db queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

// e: start a transaction => tx
// f: tx => q
// e: run the fn
// e: commit or rollback
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("tx err: %v,  rb err: %v", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId int64 `json:"to_account_id"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

// In this transaction,
// 1. create a transfer record
// 2. create entry of from_account
// 3. create entry of to_account
// 4. update from_account
// 5. update to_account

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error


		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromCcountID: arg.FromAccountId,
			ToCcountID: arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		if(arg.FromAccountId < arg.ToAccountId) {
			result.FromAccount, err = addMoney(q, arg.FromAccountId, -arg.Amount)
	
			result.ToAccount, err = addMoney(q, arg.ToAccountId, arg.Amount)
		}else {
			result.ToAccount, err = addMoney(q, arg.ToAccountId, arg.Amount)

			result.FromAccount, err = addMoney(q, arg.FromAccountId, -arg.Amount)
		}
		
		return nil
	})

	return result, err
}

func addMoney(q *Queries, accountID int64, amount int64) (account Account,err error) {
	account, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		ID: accountID,
		Amount: amount,
	})
	if err != nil {
		return 
	}

	return 
}