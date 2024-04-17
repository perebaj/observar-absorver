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

	var rEssays []marinho.RawEssay
	var err error
	if flag.Lookup("fetch").Value.String() == "true" {
		slog.Info("fetching from website")
		rEssays, err = marinho.FetchPageHTML()
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
	} else {
		slog.Info("reading from file")
		//read the raw data from the json file
		f, err := os.Open("rawessays.json")
		if err != nil {
			panic(err)
		}

		defer f.Close()

		rEssays = []marinho.RawEssay{}
		err = json.NewDecoder(f).Decode(&rEssays)
		if err != nil {
			panic(err)
		}
	}

	var parsedEssays marinho.Essays
	for _, rawEssay := range rEssays {
		slog.Info("parsing", "url", rawEssay.URL)
		parsedEssays, err = marinho.ParseHTML2Essay(rawEssay.HTML)
		if err != nil {
			slog.Error("error parsing the html", "url", rawEssay.URL, "error", err)
			continue
		}
		slog.Info("parsed", "essays", len(parsedEssays))
	}

	f, err := os.Create("essays.json")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	j, err := json.Marshal(parsedEssays)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(j)
	if err != nil {
		panic(err)
	}
}
