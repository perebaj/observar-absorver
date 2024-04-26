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
	ID        string          `db:"id" json:"id"`
	Title     string          `db:"title" json:"title"`
	URL       string          `db:"url" json:"url"`
	Content   string          `db:"content" json:"content"`
	Date      string          `db:"date" json:"date"`
	Embedding pgvector.Vector `db:"embedding" json:"embedding"`
	ModelName string          `db:"model_name" json:"model_name"`
	Dimension int             `db:"dimension" json:"dimension"`
}

type EssayResponse struct {
	ID               string  `json:"id" db:"id"`
	Title            string  `json:"title" db:"title"`
	URL              string  `json:"url" db:"url"`
	Content          string  `json:"content" db:"content"`
	Date             string  `json:"date" db:"date"`
	CosineSimilarity float64 `json:"cosine_similarity" db:"cosine_similarity"`
}

// SaveEssay saves an essay to the database
func (s *Storage) SaveEssay(essay Essay) error {
	_, err := s.db.NamedExec(
		`INSERT INTO essays (id, title, url, content, date, embedding, model_name, dimension)
		VALUES (:id, :title, :url, :content, :date, :embedding, :model_name, :dimension)
		ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		url = EXCLUDED.url,
		content = EXCLUDED.content,
		date = EXCLUDED.date,
		embedding = EXCLUDED.embedding,
		model_name = EXCLUDED.model_name,
		dimension = EXCLUDED.dimension`,
		essay,
	)
	return err
}
