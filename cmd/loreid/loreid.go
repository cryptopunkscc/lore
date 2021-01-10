package main

import (
	"encoding/hex"
	"fmt"
	"github.com/cryptopunkscc/lore/id"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "Usage: loreid <file|id>")
		os.Exit(1)
	}

	path := os.Args[1]

	if i, err := id.Parse(path); err == nil {
		fmt.Printf("%s %d\n", hex.EncodeToString(i.Hash[:]), i.Size)
		os.Exit(0)
	}

	file, err := os.Open(path)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Couldn't open file:", err)
		os.Exit(1)
	}

	resolver := id.NewResolver()
	_, err = io.Copy(resolver, file)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Couldn't open file:", err)
		os.Exit(1)
	}

	fmt.Println(resolver.Resolve())
}
