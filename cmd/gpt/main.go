// Package main in the cmd/gpt folder gather the start point of the GPT application
package main

import (
	"log/slog"

	"github.com/perebaj/marinho"
)

func main() {
	e := marinho.Essay2Embedding("Once upon a time")
	slog.Info("embedding", "embedding", e)
}
