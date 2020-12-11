package storage

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/story"
	"github.com/cryptopunkscc/lore/story/core"
	"gorm.io/gorm"
	"io/ioutil"
)

type Index struct {
	db *gorm.DB

	story.HeaderRepo
	core.FileInfoIndex
}

func NewIndex(db *gorm.DB) *Index {
	idx := &Index{
		db: db,
	}
	idx.HeaderRepo, _ = story.NewHeaderRepoGorm(db)
	idx.FileInfoIndex = core.NewFileInfoIndexGorm(db)

	return idx
}

func (idx *Index) AddData(data []byte) error {
	// Check if data contains a story
	header, err := story.ParseHeader(data)
	if err != nil {
		return err
	}

	// Resolve story ID
	storyId, err := id.ResolveID(data, nil)
	if err != nil {
		return err
	}

	// Index story header
	err = idx.HeaderRepo.Add(storyId, header)
	if err != nil {
		return err
	}

	// Extra indexes
	switch header.Type {
	case core.FileInfoStoryType:
		var s core.FileInfo
		if story.ParseStory(data, &s) == nil {
			_ = idx.FileInfoIndex.Add(storyId, s)
		}
	}

	return nil
}

func (idx *Index) AddFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return idx.AddData(data)
}

func (idx *Index) Remove(id string) error {
	_ = idx.FileInfoIndex.Remove(id)

	return idx.HeaderRepo.Remove(id)
}
