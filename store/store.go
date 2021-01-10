package store

import "github.com/cryptopunkscc/lore/id"

// Reader defines methods for reading data from a store
type Reader interface {
	Read(id id.ID) (ReadSeekCloser, error)
	List() (id.Set, error)
}

// Editor defines methods for modifying data in the store
type Editor interface {
	Create() (Writer, error)
	Delete(id id.ID) error
	Free() (uint64, error)
}

// Store includes methods for full store access
type Store interface {
	Reader
	Editor
}

// Writer defines methods for writing a new item into the store
type Writer interface {
	Write(data []byte) (int, error)
	Finalize() (id.ID, error)
	Discard() error
}

type Observer interface {
	Added(id id.ID)
	Removed(id id.ID)
}
