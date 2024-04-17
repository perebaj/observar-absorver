package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/perebaj/marinho"
)

func main() {
	rEssays, err := marinho.FetchPageHTML()
	slog.Info("fetched", "pages", len(rEssays))
	if err != nil {
		panic(err)
	}

	//save the raw result to a json file
	f, err := os.Create("rawessays.json")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	j, err := json.Marshal(rEssays)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(j)
	if err != nil {
		panic(err)
	}

	for _, rawEssay := range rEssays {
		slog.Info("parsing", "url", rawEssay.URL)
		parsedEssays, err := marinho.ParseHTML2Essay(rawEssay.HTML)
		if err != nil {
			slog.Error("error parsing the html", "url", rawEssay.URL, "error", err)
			continue
		}
		slog.Info("parsed", "essays", len(parsedEssays))
	}
	slog.Info("All essays parsed successfully!")
}
