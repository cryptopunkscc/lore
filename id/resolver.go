package id

import (
	"io"
	"os"
)

type Resolver interface {
	Write([]byte) (int, error)
	Resolve() string
}

var defaultResolver = NewID0Resolver()

func ResolveFileID(path string, resolver Resolver) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if resolver == nil {
		resolver = defaultResolver
	}

	_, err = io.Copy(resolver, file)
	if err != nil {
		return "", err
	}

	return resolver.Resolve(), nil
}
