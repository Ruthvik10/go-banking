package db

import (
	"context"
	"testing"

	"github.com/Ruthvik10/go-banking/util"
	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T, acc *Account) *Entry {
	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}
	entry, err := testQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, entry.AccountID)
	require.NotZero(t, entry.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	acc := createTestAccount(t)
	createTestEntry(t, acc)
}

func TestGetEntryByID(t *testing.T) {
	acc := createTestAccount(t)
	newEntry := createTestEntry(t, acc)

	entry, err := testQueries.GetEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.ID, newEntry.ID)
}

func TestListEntries(t *testing.T) {
	acc := createTestAccount(t)
	for i := 0; i < 10; i++ {
		createTestEntry(t, acc)
	}

	args := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    0,
	}

	accounts, err := testQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)
}
