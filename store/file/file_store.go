package file

import (
	_id "github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/util"
	"github.com/minio/minio/pkg/disk"
	"os"
	"path/filepath"
)

var _ store.Store = &FileStore{}

type FileStore struct {
	rootDir     string
	addedFunc   _id.IDFunc
	removedFunc _id.IDFunc
}

func NewFileStore(rootDir string, added _id.IDFunc, removed _id.IDFunc) (*FileStore, error) {
	store := &FileStore{
		addedFunc:   added,
		removedFunc: removed,
	}

	store.rootDir, _ = util.ExpandPath(rootDir)

	// Make sure the directory exists
	err := os.MkdirAll(store.rootDir, 0700)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (store FileStore) Free() (uint64, error) {
	info, err := disk.GetInfo(store.rootDir)
	if err != nil {
		return 0, err
	}

	return info.Free, nil
}

func (store FileStore) Read(id _id.ID) (util.ReadSeekCloser, error) {
	path := filepath.Join(store.rootDir, id.String())

	return os.OpenFile(path, os.O_RDONLY, 0)
}

func (store FileStore) List() (_id.Set, error) {
	files, err := filepath.Glob(filepath.Join(store.rootDir, "id1*"))
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

func (store FileStore) Create() (store.Writer, error) {
	writer, err := NewFileWriter(store.rootDir)
	if err != nil {
		return nil, err
	}

	// If there's no observer return the original writer directly
	if store.addedFunc == nil {
		return writer, nil
	}

	// Wrap the writer into a callback
	return NewWrappedWriter(writer, func(id _id.ID, err error) error {
		if err != nil {
			return nil
		}
		store.addedFunc(id)
		return nil
	}), nil
}

func (store FileStore) Delete(id _id.ID) error {
	path := filepath.Join(store.rootDir, id.String())

	err := os.Remove(path)
	if err != nil {
		return err
	}

	if store.removedFunc != nil {
		store.removedFunc(id)
	}

	return nil
}
