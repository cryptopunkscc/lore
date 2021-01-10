package pathcache

import "github.com/cryptopunkscc/lore/id"

type Repo interface {
	Create(Entry) error
	Delete(path string) error
	Exists(path string) bool
	All() ([]Entry, error)
	Find(path string) (Entry, error)
	FindByID(id id.ID) ([]Entry, error)
}
