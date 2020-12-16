package id

import (
	"io"
	"os"
)

// Resolver is an interface for ID resolvers
type Resolver interface {
	io.Writer
	Resolve() string
}

// DefaultResolver returns a copy of the default resolver.
func DefaultResolver() Resolver {
	return NewID1Resolver()
}

// ResolveID resolves the ID of data in a byte array. If no resolver is provided, DefaultResolver is used.
func ResolveID(data []byte, resolver Resolver) (string, error) {
	var err error

	if resolver == nil {
		resolver = DefaultResolver()
	}

	_, err = resolver.Write(data)
	if err != nil {
		return "", err
	}

	return resolver.Resolve(), nil
}

// ResolveFileID resolves the ID of file. If no resolver is provided, DefaultResolver is used.
func ResolveFileID(path string, resolver Resolver) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if resolver == nil {
		resolver = DefaultResolver()
	}

	_, err = io.Copy(resolver, file)
	if err != nil {
		return "", err
	}

	return resolver.Resolve(), nil
}
