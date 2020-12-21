package store

import (
	"errors"
	"os"
	"time"
)

var _ Store = &FileMapStore{}

// FileMapStore is an implementation of Reader that existing files to their IDs
type FileMapStore struct {
	//rootDir     string
	fileMapRepo FileMapRepo
}

// NewFileMapStore returns a new FileMapStore using provided FileMapRepo. If none is provided,
// a new FileMapRepoMem will be used.
func NewFileMapStore(fileMapRepo FileMapRepo) (*FileMapStore, error) {
	store := &FileMapStore{
		fileMapRepo: fileMapRepo,
	}

	// If no repo was provided, use a non-persistent memory-based one
	if store.fileMapRepo == nil {
		store.fileMapRepo = NewFileMapRepoMem()
	}

	return store, nil
}

// Add adds a path to id mapping to the store
func (store *FileMapStore) Add(path string, id string) error {
	return store.fileMapRepo.Create(FileMapEntry{
		Path:    path,
		ID:      id,
		AddedAt: time.Now(),
	})
}

// Remove removes a path mapping from the store
func (store *FileMapStore) Remove(path string) error {
	return store.fileMapRepo.Delete(path)
}

// Entries returns all FileMapEntry from the store
func (store *FileMapStore) Entries() ([]FileMapEntry, error) {
	return store.fileMapRepo.All()
}

func (store *FileMapStore) Contains(path string) bool {
	return store.fileMapRepo.Contains(path)
}

// Read implements Reader interface
func (store FileMapStore) Read(id string) (ReadSeekCloser, error) {
	// Get all paths from the database
	paths, err := store.fileMapRepo.Paths(id)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		file, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err == nil {
			return file, nil
		}
	}

	// TODO: Should we also automatically remove invalid entries?
	return nil, errors.New("not found")
}

// List implements Reader interface
func (store FileMapStore) List() ([]string, error) {
	var list = make([]string, 0)

	// Fetch all entries from the repo
	entries, err := store.fileMapRepo.All()
	if err != nil {
		return nil, err
	}

	// Convert to []string
	for _, entry := range entries {
		list = append(list, entry.ID)
	}

	return list, nil
}

// Create implements Editor interface. This method is unsupported for this store.
func (store FileMapStore) Create() (Writer, error) {
	return nil, ErrUnsupported
}

// Delete implements Editor interface. This method is unsupported for this store.
func (store FileMapStore) Delete(string) error {
	return ErrUnsupported
}
