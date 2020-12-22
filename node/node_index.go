package node

import (
	"github.com/cryptopunkscc/lore/index"
	"github.com/cryptopunkscc/lore/store"
	"gorm.io/gorm"
	"log"
)

type NodeIndex struct {
	fileTypeIndex    *index.FileTypeIndex
	storyHeaderIndex *index.StoryHeaderIndex
	coreLabelIndex   *index.CoreLabelIndex
}

func NewNodeIndex(db *gorm.DB) (*NodeIndex, error) {
	var err error
	var idx = &NodeIndex{}

	fileTypeRepo, err := index.NewFileTypeRepoGorm(db)
	if err != nil {
		return nil, err
	}

	idx.fileTypeIndex = index.NewFileTypeIndex(fileTypeRepo)

	storyHeaderRepo, err := index.NewStoryHeaderRepoGorm(db)
	if err != nil {
		return nil, err
	}

	idx.storyHeaderIndex = index.NewStoryHeaderIndex(storyHeaderRepo)

	coreLabelRepo, err := index.NewCoreLabelRepoGorm(db)
	if err != nil {
		return nil, err
	}

	idx.coreLabelIndex = index.NewCoreLabelIndex(coreLabelRepo)

	return idx, nil
}

func (idx *NodeIndex) Add(id string, store store.Reader) error {
	var err error

	err = idx.fileTypeIndex.Add(id, store)
	if err != nil {
		return err
	}

	err = idx.storyHeaderIndex.Add(id, store)
	if err != nil {
		return err
	}

	err = idx.coreLabelIndex.Add(id, store)
	if err != nil {
		log.Println("CoreLabelIndex:", err)
	}

	return nil
}

func (idx *NodeIndex) Remove(id string) error {
	_ = idx.fileTypeIndex.Remove(id)

	_ = idx.storyHeaderIndex.Remove(id)

	_ = idx.coreLabelIndex.Remove(id)

	return nil
}
