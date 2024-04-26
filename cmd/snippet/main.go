// Package main in the cmd/snippet folder gather the start point of the snippet application
package main

import (
	"os"

	"github.com/perebaj/marinho"
)

func main() {
	f, err := os.Open("essays.json")
	if err != nil {
		panic(err)
	}

	marinho.Snippet(f)
}
