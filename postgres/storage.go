package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	// ErrSellerNotFound is returned when the seller is not found in the database.
	ErrSellerNotFound = errors.New("seller not found")
)

// Storage deal with the database layer for transactions.
type Storage struct {
	db *sqlx.DB
}

// NewStorage initialize a new transaction storage.
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}
