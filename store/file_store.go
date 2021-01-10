package store

import (
	_id "github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/util"
	"github.com/minio/minio/pkg/disk"
	"os"
	"path/filepath"
)

var _ Store = &FileStore{}

type FileStore struct {
	rootDir      string
	addedEvent   _id.IDFunc
	removedEvent _id.IDFunc
}

func NewFileStore(rootDir string, added _id.IDFunc, removed _id.IDFunc) (*FileStore, error) {
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

func (f FileStore) Free() (uint64, error) {
	info, err := disk.GetInfo(f.rootDir)
	if err != nil {
		return 0, err
	}

	return info.Free, nil
}

func (f FileStore) Read(id _id.ID) (ReadSeekCloser, error) {
	path := filepath.Join(f.rootDir, id.String())

	return os.OpenFile(path, os.O_RDONLY, 0)
}

func (f FileStore) List() (_id.Set, error) {
	files, err := filepath.Glob(filepath.Join(f.rootDir, "id1*"))
	if err != nil {
		return nil, err
	}

	set := _id.NewSet()
	for _, file := range files {
		id, err := _id.Parse(filepath.Base(file))
		if err != nil {
			continue
		}

		set.Add(id)
	}
	return set, nil
}

func (f FileStore) Create() (Writer, error) {
	writer, err := NewFileWriter(f.rootDir)
	if err != nil {
		return nil, err
	}

	// If there's no observer return the original writer directly
	if f.addedEvent == nil {
		return writer, nil
	}

	// Wrap the writer into a callback
	return NewWrappedWriter(writer, func(id _id.ID, err error) error {
		if err != nil {
			return nil
		}
		f.addedEvent(id)
		return nil
	}), nil
}

func (f FileStore) Delete(id _id.ID) error {
	path := filepath.Join(f.rootDir, id.String())

	err := os.Remove(path)
	if err != nil {
		return err
	}

	if f.removedEvent != nil {
		f.removedEvent(id)
	}

	return nil
}
