package postgres

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
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

type Essay struct {
	Title      string          `db:"title"`
	Body       string          `db:"body"`
	BodyEmbedding  pgvector.Vector `db:"body_embedding"`
	Model_name string          `db:"model_name"`
}

// Save the essay in the database
func (s *Storage) SaveEssay(essay Essay) error {
	_, err := s.db.NamedExec(`INSERT INTO essays (title, body, body_embedding, model_name) VALUES (:title, :body, :body_embedding, :model_name)`, essay)
	return err
}
