package index

import "github.com/cryptopunkscc/lore/story"

type CoreLabelRepo interface {
	Add(id string, coreLabel story.CoreLabel) error
	Remove(id string) error
	Search(query string) ([]string, error)
}
