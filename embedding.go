// Package marinho gather all the code related to the business logic
package marinho

import (
	"context"
	"log/slog"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Config is a struct that contains the configuration to connect to the OpenAI API
type Config struct {
	OpenAIKey string
}

// EmbeddingDimension is a enum that represents possible dimensions for the embeddings
type EmbeddingDimension int

const (
	// EmbeddingDimension1536 is the dimension of the embedding. (1536 floats)
	EmbeddingDimension1536 EmbeddingDimension = 1536
)

// Essay represents all the information about an essay including the embedding
type Essay struct {
	// ID is a key formed by the title and the chunk number. ex: essay-title-0, essay-title-1
	ID      string `json:"id"`
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
	Date    string `json:"date"`
	EmbeddingEssay
}

// EmbeddingEssay is a struct that contains the embedding of an essay and some metadata about the model and the dimension that generated the embedding
type EmbeddingEssay struct {
	Text      string                `json:"text"`
	Embedding []float32             `json:"embedding"`
	ModelName openai.EmbeddingModel `json:"model_name"`
	Dimension EmbeddingDimension    `json:"dimension"`
}

// Essay2Embedding converts an essay to an embedding
func Essay2Embedding(text string) EmbeddingEssay {
	cfg := Config{
		OpenAIKey: os.Getenv("OPENAI_KEY"),
	}

	if cfg.OpenAIKey == "" {
		slog.Error("OPENAI_KEY is required")
		return EmbeddingEssay{}
	}

	client := openai.NewClient(cfg.OpenAIKey)

	queryReq := openai.EmbeddingRequest{
		Model:          openai.SmallEmbedding3,
		Input:          []string{text},
		EncodingFormat: openai.EmbeddingEncodingFormatBase64,
	}

	resp, err := client.CreateEmbeddings(context.Background(), queryReq)
	if err != nil {
		slog.Error("error embedding", "error", err)
		return EmbeddingEssay{}
	}
	e := EmbeddingEssay{
		Text:      text,
		Embedding: resp.Data[0].Embedding,
		ModelName: resp.Model,
		Dimension: EmbeddingDimension1536,
	}
	return e
}
