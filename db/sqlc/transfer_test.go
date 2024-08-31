package db

import (
	"context"
	"testing"

	"github.com/Ruthvik10/go-banking/util"
	"github.com/stretchr/testify/require"
)

func createTestTransfer(t *testing.T, fromAcc *Account, toAcc *Account) *Transfer {
	args := CreateTransferParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
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

	entry, err := testQueries.GetTransfer(context.Background(), newTransfer.ID)

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

	args := ListTransfersParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Limit:         5,
		Offset:        0,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfers)
	require.Len(t, transfers, 5)
}
