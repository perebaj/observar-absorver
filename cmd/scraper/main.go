package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/perebaj/marinho"
)

func main() {
	flag.Bool("fetch", false, "fetch the data from the website")

	flag.Parse()

	var rawEssays []marinho.RawEssay
	var err error
	if flag.Lookup("fetch").Value.String() == "true" {
		slog.Info("fetching from website")
		rawEssays, err = marinho.FetchPageHTML()
		slog.Info("fetched", "pages", len(rawEssays))
		if err != nil {
			panic(err)
		}

		f, err := os.Create("rawessays.json")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		j, err := json.Marshal(rawEssays)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(j)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Info("reading from file")

		f, err := os.Open("rawessays.json")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		rawEssays = []marinho.RawEssay{}
		err = json.NewDecoder(f).Decode(&rawEssays)
		if err != nil {
			panic(err)
		}
	}

	var parsedEssaysResult marinho.Essays
	for _, rawEssay := range rawEssays {
		slog.Info("parsing", "url", rawEssay.URL)
		parsedEssays, err := marinho.ParseHTML2Essay(rawEssay.HTML)
		if err != nil {
			slog.Error("error parsing the html", "url", rawEssay.URL, "error", err)
			continue
		}
		parsedEssaysResult = append(parsedEssaysResult, parsedEssays...)
		slog.Info("parsed", "essays", len(parsedEssays))
	}

	f, err := os.Create("essays.json")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	j, err := json.Marshal(parsedEssaysResult)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(j)
	if err != nil {
		panic(err)
	}
}
