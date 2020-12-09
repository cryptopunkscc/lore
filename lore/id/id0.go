package id

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type ID0Resolver struct {
	hash hash.Hash
}

func NewID0Resolver() *ID0Resolver {
	return &ID0Resolver{
		hash: sha256.New(),
	}
}

func (res *ID0Resolver) Write(data []byte) (int, error) {
	return res.hash.Write(data)
}

func (res *ID0Resolver) Resolve() string {
	id := "id0" + hex.EncodeToString(res.hash.Sum(nil))
	return id
}
