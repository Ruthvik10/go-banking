package store

import "errors"

var (
	ErrAccountNotFound = errors.New("AccountStore not found")
)

var (
	ErrEntryRecordNotFound = errors.New("entry not found")
)
