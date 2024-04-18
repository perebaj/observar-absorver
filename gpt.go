package marinho

import (
	"context"
	"log/slog"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type Config struct {
	OpenAIKey string
}

type EmbeddingDimension int

const (
	EmbeddingDimension_1536 EmbeddingDimension = 1536
)

type EmbeddingEssay struct {
	Embedding openai.Embedding
	Model     openai.EmbeddingModel
	Dimension EmbeddingDimension
}

func Essay2Embedding(essay string) EmbeddingEssay {
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
		Input:          []string{essay},
		EncodingFormat: openai.EmbeddingEncodingFormatBase64,
	}

	resp, err := client.CreateEmbeddings(context.Background(), queryReq)
	if err != nil {
		slog.Error("error embedding", "error", err)
		return EmbeddingEssay{}
	}

	e := EmbeddingEssay{
		Embedding: resp.Data[0],
		Model:     resp.Model,
		Dimension: EmbeddingDimension_1536,
	}

	return e
}
