package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id" db:"account_id"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type EntryStore struct {
	db DBTX
}

func NewEntryStore(db DBTX) *EntryStore {
	return &EntryStore{db}
}

func (r *EntryStore) CreateEntry(ctx context.Context, entry *Entry) error {
	q := `
			INSERT INTO entries
			(account_id, amount)
			VALUES ($1, $2)
			RETURNING *;
		`
	args := []any{entry.AccountID, entry.Amount}
	err := r.db.QueryRowxContext(ctx, q, args...).StructScan(entry)
	if err != nil {
		return err
	}
	return nil
}

func (r *EntryStore) GetEntry(ctx context.Context, id int64) (*Entry, error) {
	q := `
		SELECT * FROM entries
		WHERE id = $1
		LIMIT 1;
	`
	var entry Entry
	err := r.db.QueryRowxContext(ctx, q, id).StructScan(&entry)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrEntryRecordNotFound
		default:
			return nil, err
		}
	}
	return &entry, nil
}

type EntryArgs struct {
	AccountID int64
	Limit     int
	Offset    int
}

func (r *EntryStore) ListEntries(ctx context.Context, args *EntryArgs) ([]*Entry, error) {
	q := `
			SELECT * FROM entries
			WHERE account_id = $1
			ORDER BY id
			LIMIT $2
			OFFSET $3;
		`
	var entries []*Entry
	err := r.db.SelectContext(ctx, &entries, q, args.AccountID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
