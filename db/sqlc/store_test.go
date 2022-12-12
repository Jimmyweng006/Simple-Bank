package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/Jimmyweng006/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	fmt.Println(">> before:", fromAccount.Balance, toAccount.Balance)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	n := 5
	amount := util.RandomMoney()
	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   toAccount.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check elements in errs and results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, fromAccount.ID, transfer.FromAccountID)
		require.Equal(t, toAccount.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromAccount.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toAccount.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account's ID
		fromAccountOfResult := result.FromAccount
		require.NotEmpty(t, fromAccountOfResult)
		require.Equal(t, fromAccount.ID, fromAccountOfResult.ID)

		toAccountOfResult := result.ToAccount
		require.NotEmpty(t, toAccountOfResult)
		require.Equal(t, toAccount.ID, toAccountOfResult.ID)

		// check account's Balance
		fmt.Println(">> tx:", fromAccountOfResult.Balance, toAccountOfResult.Balance)
		diffOfFromAccount := fromAccount.Balance - fromAccountOfResult.Balance
		diffOfToAccount := toAccountOfResult.Balance - toAccount.Balance
		require.Equal(t, diffOfFromAccount, diffOfToAccount)
		require.True(t, diffOfFromAccount > 0)
		require.True(t, diffOfFromAccount%amount == 0)

		k := int(diffOfFromAccount / amount)
		require.True(t, k >= 1 && k <= n)
		// for loop element in this block are unordered
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balances
	updatedFromAccount, err := testQueries.GetAccount(context.Background(), fromAccount.ID)
	require.NoError(t, err)

	updatedToAccount, err := testQueries.GetAccount(context.Background(), toAccount.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedFromAccount.Balance, updatedToAccount.Balance)
	require.Equal(t, fromAccount.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, toAccount.Balance+int64(n)*amount, updatedToAccount.Balance)
}
