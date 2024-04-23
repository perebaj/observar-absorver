// Package postgres gather all the code related to the postgres database
package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
)

// Storage is a struct that contains the database connection
type Storage struct {
	db *sqlx.DB
}

// NewStorage creates a new storage
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Essay is a struct that represents an essay in the database
type Essay struct {
	Title         string          `db:"title"`
	Body          string          `db:"body"`
	BodyEmbedding pgvector.Vector `db:"body_embedding"`
	ModelName     string          `db:"model_name"`
}

// SaveEssay saves an essay to the database
func (s *Storage) SaveEssay(essay Essay) error {
	_, err := s.db.NamedExec(`INSERT INTO essays (title, body, body_embedding, model_name) VALUES (:title, :body, :body_embedding, :model_name)`, essay)
	return err
}
