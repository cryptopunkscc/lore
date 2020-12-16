package store

type FileMapRepo interface {
	Create(FileMapEntry) error
	Delete(path string) error
	Contains(path string) bool
	All() ([]FileMapEntry, error)
	Paths(id string) ([]string, error)
}
