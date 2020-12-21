package index

import (
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/util"
)

type FileTypeIndex struct {
	repo FileTypeRepo
}

func NewFileTypeIndex(repo FileTypeRepo) *FileTypeIndex {
	return &FileTypeIndex{
		repo: repo,
	}
}

func (idx *FileTypeIndex) Add(id string, source store.Reader) error {
	file, err := source.Read(id)
	if err != nil {
		return err
	}
	defer file.Close()

	contentType, err := util.GetContentType(file)
	if err != nil {
		return err
	}

	return idx.repo.Set(id, contentType)
}

func (idx *FileTypeIndex) Remove(id string) error {
	return idx.repo.Clear(id)
}
