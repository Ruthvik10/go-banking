package store

import (
	"context"
)

type TransferMoneyArgs struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
}

type TransferMoneyResult struct {
	FromAccount *Account  `json:"from_account"`
	ToAccount   *Account  `json:"to_account"`
	FromEntry   *Entry    `json:"from_entry"`
	ToEntry     *Entry    `json:"to_entry"`
	Transfer    *Transfer `json:"transfer"`
}

func (s *Store) TransferMoney(ctx context.Context, args *TransferMoneyArgs) (*TransferMoneyResult, error) {
	var result TransferMoneyResult
	result.Transfer = &Transfer{
		FromAccountID: args.FromAccountID,
		ToAccountID:   args.ToAccountID,
		Amount:        args.Amount,
	}
	result.FromEntry = &Entry{
		AccountID: args.FromAccountID,
		Amount:    -args.Amount,
	}
	result.ToEntry = &Entry{
		AccountID: args.ToAccountID,
		Amount:    args.Amount,
	}
	err := s.withTx(ctx, func(store *Store) error {
		var err error
		err = store.CreateTranfer(ctx, result.Transfer)
		if err != nil {
			return err
		}
		err = store.CreateEntry(ctx, result.FromEntry)
		if err != nil {
			return err
		}
		err = store.CreateEntry(ctx, result.ToEntry)
		if err != nil {
			return err
		}

		result.FromAccount, err = store.GetAccount(ctx, args.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = store.UpdateAccount(ctx, result.FromAccount.ID, result.FromAccount.Balance-args.Amount)
		if err != nil {
			return err
		}

		result.ToAccount, err = store.GetAccount(ctx, args.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = store.UpdateAccount(ctx, result.ToAccount.ID, result.ToAccount.Balance+args.Amount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
