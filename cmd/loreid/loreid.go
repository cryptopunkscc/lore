package main

import (
	"fmt"
	"github.com/cryptopunkscc/lore/id"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: loreid <file>")
		os.Exit(1)
	}

	path := os.Args[1]

	file, err := os.Open(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Couldn't open file:", err)
		os.Exit(1)
	}

	resolver := id.NewID1Resolver()
	_, err = io.Copy(resolver, file)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Couldn't open file:", err)
		os.Exit(1)
	}

	fmt.Println(resolver.Resolve())
}
