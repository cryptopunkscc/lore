package store

import (
	"github.com/cryptopunkscc/lore/util"
	"github.com/minio/minio/pkg/disk"
	"os"
	"path/filepath"
)

var _ Store = &FileStore{}

type FileStore struct {
	rootDir string
}

func (f FileStore) Free() (int64, error) {
	info, err := disk.GetInfo(f.rootDir)
	if err != nil {
		return 0, err
	}

	return int64(info.Free), nil
}

func NewFileStore(rootDir string) (*FileStore, error) {
	store := &FileStore{}

	store.rootDir, _ = util.ExpandPath(rootDir)

	// Make sure the directory exists
	err := os.MkdirAll(store.rootDir, 0700)
	if err != nil {
		return nil, err
	}

	return store, nil
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
	return NewFileWriter(f.rootDir, nil)
}

func (f FileStore) Delete(id string) error {
	path := filepath.Join(f.rootDir, id)

	return os.Remove(path)
}
