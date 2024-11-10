package store

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferMoney(t *testing.T) {
	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)
	errs := make(chan error)
	results := make(chan *TransferMoneyResult)

	n := 5
	amount := int64(10)

	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferMoney(context.Background(), &TransferMoneyArgs{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	// check results
	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acc1.ID, transfer.FromAccountID)
		require.Equal(t, acc2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acc1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, acc2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, toAccount.ID)

		// check balances
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)

		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updatedfromAccount, err := store.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)

	updatedtoAccount, err := store.GetAccount(context.Background(), acc2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedfromAccount.Balance, updatedtoAccount.Balance)

	require.Equal(t, acc1.Balance-int64(n)*amount, updatedfromAccount.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedtoAccount.Balance)
}
