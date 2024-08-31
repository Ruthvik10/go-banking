package db

import (
	"context"
	"testing"

	"github.com/Ruthvik10/go-banking/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) *Account {
	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	acc, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	require.Equal(t, args.Owner, acc.Owner)
	require.Equal(t, args.Balance, acc.Balance)
	require.Equal(t, args.Currency, acc.Currency)
	return acc
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	newAcc := createTestAccount(t)
	acc, err := testQueries.GetAccount(context.Background(), newAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, newAcc.ID, acc.ID)
	require.Equal(t, newAcc.Owner, acc.Owner)
}

func TestUpdateAccount(t *testing.T) {
	newAcc := createTestAccount(t)

	args := UpdateAccountParams{
		ID:      newAcc.ID,
		Balance: util.RandomMoney(),
	}

	_, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(context.Background(), args.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, newAcc.ID, acc.ID)
	require.NotEqual(t, newAcc.Balance, acc.Balance)
}

func TestDeleteAccount(t *testing.T) {
	newAcc := createTestAccount(t)

	err := testQueries.DeleteAccount(context.Background(), newAcc.ID)
	require.NoError(t, err)

	acc, err := testQueries.GetAccount(context.Background(), newAcc.ID)
	require.Error(t, err)
	require.Empty(t, acc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, a := range accounts {
		require.NotEmpty(t, a)
	}
}
