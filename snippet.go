// Package marinho snippet.go receive a list of essays and chunk them into smaller pieces.
package marinho

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

// EssayChunk is a struct that gather the information of a chunk of an essay, splitted by paragraphs.
type EssayChunk struct {
	ID      string
	Title   string
	Content string
	Date    string
	URL     string
}

// Snippet receives a file with a list of essays and chunk them into smaller pieces.
func Snippet(f *os.File) {
	var e ScraperEssays
	err := json.NewDecoder(f).Decode(&e)
	if err != nil {
		panic(err)
	}

	var result []EssayChunk
	for _, v := range e {
		//split the text using .
		var chunkID int
		for _, s := range strings.Split(v.Content, ".") {
			chunkIDString := strconv.Itoa(chunkID)
			var id string
			if v.Title == "" {
				uid := uuid.NewString()
				id = slug.Make(uid) + "-" + chunkIDString
			} else {
				id = slug.Make(v.Title) + "-" + chunkIDString
			}

			eChunk := EssayChunk{
				ID:      id,
				Title:   v.Title,
				Content: s,
				Date:    v.Date,
				URL:     v.URL,
			}
			chunkID++
			result = append(result, eChunk)
		}
	}

	f, err = os.Create("essays_chunk.json")
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(f).Encode(result)
	if err != nil {
		panic(err)
	}
	fmt.Println("Essays chunked successfully")
}
