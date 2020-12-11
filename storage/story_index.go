package storage

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/story"
)

type StoryIndex struct {
	storyRepo story.StoryRepo
}

func (idx *StoryIndex) IndexAs(id string, header *story.Header) error {
	var err error

	err = idx.storyRepo.SetStoryType(id, header.Type)
	if err != nil {
		return err
	}

	err = idx.storyRepo.SetStoryRels(id, header.Rel)
	if err != nil {
		return err
	}

	return nil
}

func (idx *StoryIndex) IndexFile(path string) error {
	header, err := story.ParseHeaderFromFile(path)
	if err == nil {
		fileId, err := id.ResolveFileID(path, nil)
		if err != nil {
			return err
		}
		return idx.IndexAs(fileId, header)
	}
	return nil
}

func (idx *StoryIndex) Forget(id string) error {
	return idx.storyRepo.Forget(id)
}
