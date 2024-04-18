package main

import (
	"log/slog"

	"github.com/perebaj/marinho"
)

func main() {
	e := marinho.Essay2Embedding("Once upon a time")
	slog.Info("embedding", "embedding", e)
}
