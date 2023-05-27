package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func(i int) {

			ctx := context.WithValue(context.Background(), "txID", i)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}(i)
	}

	for i := 0; i < n; i++ {
		err := <-errs
		result := <-results

		require.NoError(t, err)
		require.NotEmpty(t, result)

		// Check transfer
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, account1.ID, result.Transfer.FromAccountID)
		require.Equal(t, account2.ID, result.Transfer.ToAccountID)
		require.Equal(t, amount, result.Transfer.Amount)
		require.NotZero(t, result.Transfer.ID)
		require.NotZero(t, result.Transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		// Check From Entry
		require.NotEmpty(t, result.FromEntry)
		require.Equal(t, account1.ID, result.FromEntry.AccountID)
		require.Equal(t, -amount, result.FromEntry.Amount)
		require.NotZero(t, result.FromEntry.ID)
		require.NotZero(t, result.FromEntry.CreatedAt)

		// Check To Entry
		require.NotEmpty(t, result.ToEntry)
		require.Equal(t, account2.ID, result.ToEntry.AccountID)
		require.Equal(t, amount, result.ToEntry.Amount)
		require.NotZero(t, result.ToEntry.ID)
		require.NotZero(t, result.ToEntry.CreatedAt)

		// TODO: Check balances
		require.NotEmpty(t, result.FromAccount)
		require.NotEmpty(t, result.ToAccount)

		diff1 := account1.Balance - result.FromAccount.Balance
		diff2 := result.ToAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
	}

	fromAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fromAccount1)
	require.Equal(t, account1.Balance-(int64(n)*amount), fromAccount1.Balance)

	fromAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fromAccount2)
	require.Equal(t, account2.Balance+(int64(n)*amount), fromAccount2.Balance)

}
