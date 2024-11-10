package store

import (
	"bank-app/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) *Account {
	acc := &Account{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	err := store.CreateAccount(context.Background(), acc)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	require.NotZero(t, acc.Balance)
	require.NotZero(t, acc.Owner)
	require.NotZero(t, acc.Currency)

	return acc
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccountByID(t *testing.T) {
	newAcc := createTestAccount(t)

	acc, err := store.GetAccount(context.Background(), newAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, acc.ID, newAcc.ID)
}

func TestUpdateAccount(t *testing.T) {
	newAcc := createTestAccount(t)
	newBalance := util.RandomMoney()
	acc, err := store.UpdateAccount(context.Background(), newAcc.ID, newBalance)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, newAcc.ID, acc.ID)
	require.Equal(t, newBalance, acc.Balance)
}

func TestDeleteAccount(t *testing.T) {
	newAcc := createTestAccount(t)
	err := store.DeleteAccount(context.Background(), newAcc.ID)
	require.NoError(t, err)

	acc, err := store.GetAccount(context.Background(), newAcc.ID)
	require.Equal(t, ErrAccountNotFound, err)
	require.Empty(t, acc)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createTestAccount(t)
	}

	opts := ListAccountArgs{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := store.ListAccounts(context.Background(), &opts)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)
}
