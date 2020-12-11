package index

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/story"
	"gorm.io/gorm"
	"io/ioutil"
)

type StoryIndex struct {
	db        *gorm.DB
	storyRepo story.StoryRepo
	indexers  map[string]Indexer
}

type Indexer interface {
	Index(id string, header *story.Header, data []byte) error
	Query(query string) ([]string, error)
	Deindex(id string) error
}

func NewStoryIndex(db *gorm.DB) *StoryIndex {
	idx := &StoryIndex{
		db:       db,
		indexers: make(map[string]Indexer, 0),
	}
	idx.storyRepo, _ = story.NewStoryRepoGorm(db)

	idx.RegisterIndexer("core.fileinfo", newCoreFileinfoIndexer(idx.db))

	return idx
}

func (idx *StoryIndex) IndexFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return idx.Index(data)
}

func (idx *StoryIndex) Index(data []byte) error {
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
	err = idx.indexHeader(storyId, header)
	if err != nil {
		return err
	}

	// See if we have an indexer for this story
	indexer, ok := idx.indexers[header.Type]
	if ok {
		indexer.Index(storyId, header, data)
	}

	return nil
}

func (idx *StoryIndex) QueryType(typ string, query string) ([]string, error) {
	if i, ok := idx.indexers[typ]; ok {
		return i.Query(query)
	}

	return []string{}, nil
}

func (idx *StoryIndex) Forget(id string) error {
	for _, i := range idx.indexers {
		_ = i.Deindex(id)
	}

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

func (idx *StoryIndex) RegisterIndexer(typ string, indexer Indexer) {
	idx.indexers[typ] = indexer
}
