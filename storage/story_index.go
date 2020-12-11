package storage

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/story"
	"io/ioutil"
)

type StoryIndex struct {
	storyRepo story.StoryRepo
}

func (idx *StoryIndex) IndexFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return idx.Index(data)
}

func (idx *StoryIndex) Index(data []byte) error {
	header, err := story.ParseHeader(data)
	if err == nil {
		fileId, err := id.ResolveID(data, nil)
		if err != nil {
			return err
		}
		return idx.indexHeader(fileId, header)
	}
	return nil
}

func (idx *StoryIndex) Forget(id string) error {
	return idx.storyRepo.Forget(id)
}

func (idx *StoryIndex) indexHeader(id string, header *story.Header) error {
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
