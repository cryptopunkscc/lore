package id

import (
	"io"
	"os"
)

type Resolver interface {
	Write([]byte) (int, error)
	Resolve() string
}

func defaultResolver() Resolver {
	return NewID0Resolver()
}

func ResolveID(data []byte, resolver Resolver) (string, error) {
	var err error

	if resolver == nil {
		resolver = defaultResolver()
	}

	_, err = resolver.Write(data)
	if err != nil {
		return "", err
	}

	return resolver.Resolve(), nil
}

func ResolveFileID(path string, resolver Resolver) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if resolver == nil {
		resolver = defaultResolver()
	}

	_, err = io.Copy(resolver, file)
	if err != nil {
		return "", err
	}

	return resolver.Resolve(), nil
}
