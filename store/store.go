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
}

// ReadEditor includes methods for full store access
type ReadEditor interface {
	Reader
	Editor
}

// Writer defines methods for writing a new item into the store
type Writer interface {
	Write(data []byte) (int, error)
	Finalize() (string, error)
	Discard() error
}

// Info is a struct containing information about a store item
type Info struct {
	Size int64 // Size in bytes of the item
}
