// Package main in the cmd/embedding...
package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/perebaj/marinho"
)

func main() {
	textFlag := flag.String("text", "", "text to convert to an embedding")
	flag.Parse()

	if *textFlag == "" {
		flag.Usage()
		return
	}

	e := marinho.Essay2Embedding(*textFlag)
	eByte, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(eByte))
}
