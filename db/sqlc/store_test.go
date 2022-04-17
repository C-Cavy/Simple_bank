package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println("before: ", account1.Balance, account2.Balance)

	n := 3
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId: account2.ID,
				Amount: amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

		// check result
		result := <- results
		require.NotEmpty(t, result)

		// check result.transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromCcountID)
		require.Equal(t, account2.ID, transfer.ToCcountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		
		// check result.FromEntry
		from_entry := result.FromEntry
		require.NotEmpty(t, from_entry)
		require.NotZero(t, from_entry.ID)
		require.Equal(t, account1.ID, from_entry.AccountID)
		require.Equal(t, -amount, from_entry.Amount)
		require.NotZero(t, from_entry.CreatedAt)

		// check result.ToEntry
		to_entry := result.ToEntry
		require.NotEmpty(t, to_entry)
		require.NotZero(t, to_entry.ID)
		require.Equal(t, account2.ID, to_entry.AccountID)
		require.Equal(t, amount, to_entry.Amount)
		require.NotZero(t, to_entry.CreatedAt)

		// TODO: check accounts's balance
		from_account := result.FromAccount
		require.NotEmpty(t, from_account)
		require.Equal(t, account1.ID, from_account.ID)

		to_account := result.ToAccount
		require.NotEmpty(t, to_account)
		require.Equal(t, account2.ID, to_account.ID)

		fmt.Println("tx: ", from_account.Balance, to_account.Balance)

		diff1 := account1.Balance - from_account.Balance
		diff2 := to_account.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

		
	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, amount * int64(n) , account1.Balance - updatedAccount1.Balance)
	require.Equal(t, amount * int64(n), updatedAccount2.Balance - account2.Balance)

	fmt.Println("end: ", updatedAccount1.Balance, updatedAccount2.Balance)
}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println("before: ", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i % 2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}
		go func() {
			ctx := context.Background()
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId: toAccountID,
				Amount: amount,
			})

			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <- errs
		require.NoError(t, err)

	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, updatedAccount2.Balance, account2.Balance)

	fmt.Println("end: ", updatedAccount1.Balance, updatedAccount2.Balance)
}