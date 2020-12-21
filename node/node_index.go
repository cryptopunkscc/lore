package node

import (
	"github.com/cryptopunkscc/lore/index"
	"github.com/cryptopunkscc/lore/store"
	"gorm.io/gorm"
)

type NodeIndex struct {
	fileTypeIdx    *index.FileTypeIndex
	storyHeaderIdx *index.StoryHeaderIndex
}

func NewNodeIndex(db *gorm.DB) (*NodeIndex, error) {
	var err error
	var idx = &NodeIndex{}

	fileTypeRepo, err := index.NewFileTypeRepoGorm(db)
	if err != nil {
		return nil, err
	}

	idx.fileTypeIdx = index.NewFileTypeIndex(fileTypeRepo)

	storyHeaderRepo, err := index.NewStoryHeaderRepoGorm(db)
	if err != nil {
		return nil, err
	}

	idx.storyHeaderIdx = index.NewStoryHeaderIndex(storyHeaderRepo)

	return idx, nil
}

func (idx *NodeIndex) Add(id string, store store.Reader) error {
	var err error

	err = idx.fileTypeIdx.Add(id, store)
	if err != nil {
		return err
	}

	err = idx.storyHeaderIdx.Add(id, store)
	if err != nil {
		return err
	}

	return nil
}

func (idx *NodeIndex) Remove(id string) error {
	_ = idx.fileTypeIdx.Remove(id)

	_ = idx.storyHeaderIdx.Remove(id)

	return nil
}
