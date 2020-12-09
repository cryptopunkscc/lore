package id

import (
	"io"
	"os"
)

type Resolver interface {
	Write([]byte) (int, error)
	Resolve() string
}

func ResolveFileID(path string, resolver Resolver) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(resolver, file)
	if err != nil {
		return "", err
	}

	return resolver.Resolve(), nil
}
