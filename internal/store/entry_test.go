package store

import (
	"bank-app/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestEntry(t *testing.T, acc *Account) *Entry {
	entry := &Entry{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}
	err := store.CreateEntry(context.Background(), entry)
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

	entry, err := store.GetEntry(context.Background(), newEntry.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.ID, newEntry.ID)
}

func TestListEntries(t *testing.T) {
	acc := createTestAccount(t)
	for i := 0; i < 10; i++ {
		createTestEntry(t, acc)
	}

	opts := EntryArgs{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    0,
	}

	accounts, err := store.ListEntries(context.Background(), &opts)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)
}
