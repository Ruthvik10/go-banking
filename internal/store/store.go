package store

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	*AccountStore
	*EntryStore
	*TransferStore
	db DBTX
}

func NewStore(db DBTX) *Store {
	return &Store{
		NewAccountStore(db),
		NewEntryStore(db),
		NewTransferStore(db),
		db,
	}
}

func (s *Store) withTx(ctx context.Context, f func(store *Store) error) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	txs := NewStore(tx)
	if err := f(txs); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("rbErr: %w, rollback transaction error: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
