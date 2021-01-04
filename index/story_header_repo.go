package index

import "github.com/cryptopunkscc/lore/story"

type StoryHeaderRepo interface {
	Add(id string, header *story.Header) error
	Remove(id string) error
	Find(types []string, rels []string) ([]string, error)
}
