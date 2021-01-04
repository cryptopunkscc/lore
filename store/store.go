package store

// Reader defines methods for reading data from a store
type Reader interface {
	Read(id string) (ReadSeekCloser, error)
	List() ([]string, error)
}

// Editor defines methods for modifying data in the store
type Editor interface {
	Create() (Writer, error)
	Delete(id string) error
	Free() (int64, error)
}

// Store includes methods for full store access
type Store interface {
	Reader
	Editor
}

// Writer defines methods for writing a new item into the store
type Writer interface {
	Write(data []byte) (int, error)
	Finalize() (string, error)
	Discard() error
}
