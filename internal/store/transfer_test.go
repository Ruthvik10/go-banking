package store

import (
	"bank-app/util"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T, fromAcc *Account, toAcc *Account) *Transfer {
	transfer := &Transfer{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	}
	err := store.CreateTranfer(context.Background(), transfer)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.FromAccountID)
	require.NotZero(t, transfer.ToAccountID)
	require.NotZero(t, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAcc := createTestAccount(t)
	toAcc := createTestAccount(t)
	createTestTransfer(t, fromAcc, toAcc)
}

func TestGetTransferByID(t *testing.T) {
	fromAcc := createTestAccount(t)
	toAcc := createTestAccount(t)
	newTransfer := createTestTransfer(t, fromAcc, toAcc)

	entry, err := store.GetTransfer(context.Background(), newTransfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, entry.ID, newTransfer.ID)
}

func TestListTransfers(t *testing.T) {
	fromAcc := createTestAccount(t)
	toAcc := createTestAccount(t)
	for i := 0; i < 10; i++ {
		createTestTransfer(t, fromAcc, toAcc)
	}

	opts := TransferArgs{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := store.ListTransfers(context.Background(), &opts)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)
}
