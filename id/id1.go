package id

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/binary"
	"hash"
	"strings"
)

const id1Prefix = "id1"
const zBase32CharSet = "ybndrfg8ejkmcpqxot1uwisza345h769"
const zBase32Zero = 'y'
const id1MaxLen = 67 // Maximum length of an id1 (including prefix)

// Remove leading zeros when resolving the ID
const id1CompressionEnabled = true

// ID1 of an empty file
const ID1NullFile = "id1ba7yatbjt90hn1pxz7geufz51jb8i306e3r51pgkjfc3dphffqni"

var zBase32Encoding = base32.NewEncoding(zBase32CharSet)

// ID1Resolver is an IDResolver for ID1 format
type ID1Resolver struct {
	hash hash.Hash
	size int64
}

// NewID1Resolver returns a new instance of the ID1Resolver
func NewID1Resolver() *ID1Resolver {
	return &ID1Resolver{
		hash: sha256.New(),
	}
}

// Write more data to the buffer for ID calculation
func (res *ID1Resolver) Write(data []byte) (int, error) {
	n, err := res.hash.Write(data)
	res.size = res.size + int64(n)
	return n, err
}

// Resolve returns an ID1 of bytes written so far. You can
func (res *ID1Resolver) Resolve() string {
	var buf [40]byte
	var sum = res.hash.Sum(nil)

	// Put size and sha256 checksum in the buffer
	binary.BigEndian.PutUint64(buf[0:8], uint64(res.size))
	copy(buf[8:40], sum[0:32])

	// Encode it with zBase32
	enc := zBase32Encoding.EncodeToString(buf[:])
	if id1CompressionEnabled {
		enc = strings.TrimLeft(enc, string([]byte{zBase32Zero}))
	}

	return id1Prefix + enc
}

// ParseID1 extracts file size and sha256 checksum from an ID1 string
func ParseID1(id string) (size uint64, sha256 []byte) {
	// Check the prefix
	if !strings.HasPrefix(id, id1Prefix) {
		return 0, nil
	}

	// Check how many zeros are missing
	c := id1MaxLen - len(id)

	// Prepare the buffer with 12 leading zeros (max that compression can remove)
	var enc = [64]byte{
		zBase32Zero, zBase32Zero, zBase32Zero, zBase32Zero,
		zBase32Zero, zBase32Zero, zBase32Zero, zBase32Zero,
		zBase32Zero, zBase32Zero, zBase32Zero, zBase32Zero,
	}

	copy(enc[c:], id[3:])

	data, err := zBase32Encoding.DecodeString(string(enc[:]))
	if err != nil {
		panic(err)
	}

	size = binary.BigEndian.Uint64(data[0:8])
	sha256 = data[8:40]

	return
}

// IsID1 returns true if the string is a valid ID1, false otherwise.
func IsID1(id string) bool {
	return strings.HasPrefix(id, id1Prefix)
}
