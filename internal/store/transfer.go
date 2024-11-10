package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id" db:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id" db:"to_account_id"`
	// must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type TransferStore struct {
	db DBTX
}

func NewTransferStore(db DBTX) *TransferStore {
	return &TransferStore{db: db}
}

var (
	errTransferRecordNotFound = errors.New("entry not found")
)

func (r *TransferStore) CreateTranfer(ctx context.Context, transfer *Transfer) error {
	q := `
			INSERT INTO transfers (
			from_account_id,
			to_account_id,
			amount
			) VALUES (
			$1, $2, $3
			) RETURNING *;
		`
	args := []any{transfer.FromAccountID, transfer.ToAccountID, transfer.Amount}
	err := r.db.QueryRowxContext(ctx, q, args...).StructScan(transfer)
	if err != nil {
		return err
	}
	return nil
}

func (r *TransferStore) GetTransfer(ctx context.Context, id int64) (*Transfer, error) {
	q := `
		SELECT * FROM transfers
		WHERE id = $1
		LIMIT 1;
	`
	var entry Transfer
	err := r.db.QueryRowxContext(ctx, q, id).StructScan(&entry)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errTransferRecordNotFound
		default:
			return nil, err
		}
	}
	return &entry, nil
}

type TransferArgs struct {
	FromAccountID int64
	ToAccountID   int64
	Limit         int
	Offset        int
}

func (r *TransferStore) ListTransfers(ctx context.Context, args *TransferArgs) ([]*Transfer, error) {
	q := `
			SELECT * FROM transfers
			WHERE
				from_account_id = $1 OR
				to_account_id = $2
			ORDER BY id
			LIMIT $3
			OFFSET $4;
		`
	var transfers []*Transfer
	err := r.db.SelectContext(ctx, &transfers, q, args.FromAccountID, args.ToAccountID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	return transfers, nil
}
