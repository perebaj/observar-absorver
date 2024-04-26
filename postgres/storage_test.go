package postgres_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/perebaj/marinho"
	"github.com/perebaj/marinho/postgres"
	"github.com/pgvector/pgvector-go"
)

// OpenDB create a new database for testing and return a connection to it.
func OpenDB(t *testing.T) *sqlx.DB {
	t.Helper()

	cfg := postgres.Config{
		URL:             os.Getenv("POSTGRES_URL"),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		t.Fatalf("error connecting to Postgres: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		t.Fatalf("error pinging postgres: %v", err)
	}

	// create a new database with random suffix
	postgresURL, err := url.Parse(cfg.URL)
	if err != nil {
		t.Fatalf("error parsing Postgres connection URL: %v", err)
	}
	database := strings.TrimLeft(postgresURL.Path, "/")

	randSuffix := fmt.Sprintf("%x", time.Now().UnixNano())

	database = fmt.Sprintf("%s-%x", database, randSuffix)
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, database))
	if err != nil {
		t.Fatalf("error creating database for test: %v", err)
	}

	postgresURL.Path = "/" + database
	cfg.URL = postgresURL.String()
	testDB, err := postgres.OpenDB(cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// after run the tests, drop the database
	t.Cleanup(func() {
		defer func() {
			_ = testDB.Close()
		}()

		defer func() {
			_ = db.Close()
		}()
		_, err = db.Exec(fmt.Sprintf(`DROP DATABASE "%s" WITH (FORCE);`, database))
		if err != nil {
			t.Fatalf("error dropping database for test: %v", err)
		}
	})

	return testDB
}

func TestNewStorage(t *testing.T) {
	db := OpenDB(t)

	storage := postgres.NewStorage(db)
	if storage == nil {
		t.Fatal("storage is nil")
	}

	e := make([]float32, 1536)
	err := storage.SaveEssay(postgres.Essay{
		ID:        "1",
		Title:     "title",
		URL:       "url",
		Content:   "content",
		Date:      "date",
		Embedding: pgvector.NewVector(e),
		ModelName: "model",
		Dimension: 1536,
	})

	if err != nil {
		t.Fatalf("error saving essay: %v", err)
	}

	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM essays")

	if err != nil {
		t.Fatalf("error counting essays: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 essay, got %d", count)
	}
}

func TestSaveEssay(t *testing.T) {
	db := OpenDB(t)

	storage := postgres.NewStorage(db)
	if storage == nil {
		t.Fatal("storage is nil")
	}

	// load embeddings from a file
	f, err := os.Open("embeddingessays.json")
	if err != nil {
		t.Fatalf("error opening embeddings file: %v", err)
	}

	defer func() {
		_ = f.Close()
	}()

	var e []marinho.Essay
	err = json.NewDecoder(f).Decode(&e)
	if err != nil {
		t.Fatalf("error decoding embeddings: %v", err)
	}

	for _, v := range e {
		err = storage.SaveEssay(postgres.Essay{
			ID:        v.ID,
			Title:     v.Title,
			URL:       v.URL,
			Content:   v.Content,
			Date:      v.Date,
			Embedding: pgvector.NewVector(v.Embedding),
			ModelName: string(v.ModelName),
			Dimension: int(v.Dimension),
		})

		if err != nil {
			t.Fatalf("error saving essay: %v", err)
		}
	}

	var got int
	err = db.Get(&got, "SELECT COUNT(*) FROM essays")

	if err != nil {
		t.Fatalf("error counting essays: %v", err)
	}

	if got != len(e) {
		t.Fatalf("expected %d essays, got %d", len(e), got)
	}

	f, err = os.Open("query_input.json")
	if err != nil {
		t.Fatalf("error opening query_input file: %v", err)
	}
	var embeddingEssay marinho.EmbeddingEssay
	err = json.NewDecoder(f).Decode(&embeddingEssay)
	if err != nil {
		t.Fatalf("error decoding query_input: %v", err)
	}
	var resp []postgres.EssayResponse
	err = db.Select(&resp, "SELECT id, title, url, content, date, 1 - (embedding <=> $1) as cosine_similarity FROM essays ORDER BY cosine_similarity DESC", pgvector.NewVector(embeddingEssay.Embedding))
	if err != nil {
		t.Fatalf("error selecting essays: %v", err)
	}

	for _, v := range resp {
		t.Log(v.CosineSimilarity)
	}
}
