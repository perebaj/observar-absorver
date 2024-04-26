// Package main in the cmd/gpt folder gather the start point of the gpt application
package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/perebaj/marinho"
)

func main() {
	f, err := os.Open("essays_chunk.json")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = f.Close()
	}()

	var e []marinho.EssayChunk
	err = json.NewDecoder(f).Decode(&e)
	if err != nil {
		panic(err)
	}

	var res []marinho.Essay
	for _, v := range e[:5] {
		if v.Title != "" && v.Content != "" && v.URL != "" {
			slog.Info("creating embedding", "title", v.Title)
			e := marinho.Essay2Embedding(v.Content + v.Title)
			res = append(res, marinho.Essay{
				ID:      v.ID,
				Title:   v.Title,
				URL:     v.URL,
				Content: v.Content,
				Date:    v.Date,
				EmbeddingEssay: marinho.EmbeddingEssay{
					Embedding: e.Embedding,
					ModelName: e.ModelName,
					Dimension: marinho.EmbeddingDimension1536,
				},
			})
		}
	}
	slog.Info("saving essays", "count", len(res), "example", res[0])
	f, err = os.Create("embeddingessays.json")
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = f.Close()
	}()

	j, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(j)
	if err != nil {
		panic(err)
	}
}
