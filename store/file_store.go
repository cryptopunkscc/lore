package store

import (
	"github.com/cryptopunkscc/lore/util"
	"github.com/minio/minio/pkg/disk"
	"os"
	"path/filepath"
)

var _ Store = &FileStore{}

type FileStore struct {
	rootDir      string
	addedEvent   EventFunc
	removedEvent EventFunc
}

func NewFileStore(rootDir string, added EventFunc, removed EventFunc) (*FileStore, error) {
	store := &FileStore{
		addedEvent:   added,
		removedEvent: removed,
	}

	store.rootDir, _ = util.ExpandPath(rootDir)

	// Make sure the directory exists
	err := os.MkdirAll(store.rootDir, 0700)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (f FileStore) Free() (int64, error) {
	info, err := disk.GetInfo(f.rootDir)
	if err != nil {
		return 0, err
	}

	return int64(info.Free), nil
}

func (f FileStore) Read(id string) (ReadSeekCloser, error) {
	path := filepath.Join(f.rootDir, id)

	return os.OpenFile(path, os.O_RDONLY, 0)
}

func (f FileStore) List() ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(f.rootDir, "id*"))
	if err != nil {
		return nil, err
	}
	list := make([]string, 0)
	for _, m := range matches {
		list = append(list, filepath.Base(m))
	}
	return list, nil
}

func (f FileStore) Create() (Writer, error) {
	writer, err := NewFileWriter(f.rootDir, nil)
	if err != nil {
		return nil, err
	}

	// If there's no observer return the original writer directly
	if f.addedEvent == nil {
		return writer, nil
	}

	// Wrap the writer into a callback
	return NewWrappedWriter(writer, func(id string, err error) error {
		if err != nil {
			return nil
		}
		f.addedEvent(id)
		return nil
	}), nil
}

func (f FileStore) Delete(id string) error {
	path := filepath.Join(f.rootDir, id)

	err := os.Remove(path)
	if err != nil {
		return err
	}

	if f.removedEvent != nil {
		f.removedEvent(id)
	}

	return nil
}
