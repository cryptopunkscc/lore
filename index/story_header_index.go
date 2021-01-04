package index

import (
	"fmt"
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/story"
)

type StoryHeaderIndex struct {
	repo StoryHeaderRepo
}

func NewStoryHeaderIndex(repo StoryHeaderRepo) *StoryHeaderIndex {
	return &StoryHeaderIndex{repo: repo}
}

func (idx StoryHeaderIndex) Add(fileId string, store store.Reader) error {
	size, _ := id.ParseID1(fileId)

	// Check size from the ID to avoid pinging the store
	if size > story.MaxStorySize {
		return fmt.Errorf("not a story: %w", story.ErrDataTooBig)
	}

	// Read the file
	reader, err := store.Read(fileId)
	if err != nil {
		return err
	}

	// Get the story header
	header, err := story.ParseHeader(reader)
	if err != nil {
		return story.ErrInvalidStory
	}

	return idx.repo.Add(fileId, header)
}

func (idx StoryHeaderIndex) Remove(fileId string) error {
	return idx.repo.Remove(fileId)
}

func (idx StoryHeaderIndex) Query(types []string, rels []string) ([]string, error) {
	return nil, nil
}
