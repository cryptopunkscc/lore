package id

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"hash"
)

const id0Prefix = "id0"

type ID0Resolver struct {
	hash hash.Hash
	size int64
}

func NewID0Resolver() *ID0Resolver {
	return &ID0Resolver{
		hash: sha256.New(),
	}
}

func (res *ID0Resolver) Write(data []byte) (int, error) {
	n, err := res.hash.Write(data)
	res.size = res.size + int64(n)

	return n, err
}

func (res *ID0Resolver) Resolve() string {
	var b [8]byte
	binary.BigEndian.PutUint64(b[0:8], uint64(res.size))
	var s = hex.EncodeToString(b[0:8])
	var h = hex.EncodeToString(res.hash.Sum(nil))

	return id0Prefix + s + h
}
