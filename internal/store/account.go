package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type AccountStore struct {
	db DBTX
}

type Account struct {
	ID        int64     `json:"id" db:"id"`
	Owner     string    `json:"owner" db:"owner"`
	Balance   int64     `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func NewAccountStore(db DBTX) *AccountStore {
	return &AccountStore{db}
}

func (r *AccountStore) CreateAccount(ctx context.Context, acc *Account) error {
	q := `
			INSERT INTO accounts
			(owner, balance, currency)
			VALUES
			($1, $2, $3)
			RETURNING *;
		`
	args := []any{acc.Owner, acc.Balance, acc.Currency}
	err := r.db.QueryRowxContext(ctx, q, args...).StructScan(acc)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccountStore) GetAccount(ctx context.Context, id int64) (*Account, error) {
	q := `
		SELECT * FROM accounts
		WHERE id = $1
		FOR NO KEY UPDATE
		LIMIT 1;
	`
	var acc Account
	err := r.db.QueryRowxContext(ctx, q, id).StructScan(&acc)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrAccountNotFound
		default:
			return nil, err
		}
	}
	return &acc, nil
}

func (r *AccountStore) UpdateAccount(ctx context.Context, id, balance int64) (*Account, error) {
	q := `
			UPDATE accounts
			SET balance = $1
			WHERE id = $2
			RETURNING *;
		`
	var acc Account
	err := r.db.QueryRowxContext(ctx, q, balance, id).StructScan(&acc)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrAccountNotFound
		default:
			return nil, err
		}
	}
	return &acc, nil
}

func (r *AccountStore) DeleteAccount(ctx context.Context, id int64) error {
	q := `
			DELETE from accounts
			WHERE id = $1;
		`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrAccountNotFound
		default:
			return err
		}
	}
	return nil
}

type ListAccountArgs struct {
	Limit  int
	Offset int
}

func (r *AccountStore) ListAccounts(ctx context.Context, args *ListAccountArgs) ([]*Account, error) {
	q := `
			SELECT * FROM accounts
			ORDER BY id
			LIMIT $1
			OFFSET $2;
		`
	var accounts []*Account
	err := r.db.SelectContext(ctx, &accounts, q, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
