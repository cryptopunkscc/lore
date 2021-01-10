package pathcache

import (
	"errors"
	_id "github.com/cryptopunkscc/lore/id"
)

var _ Repo = &RepoMem{}

type RepoMem struct {
	fileMap map[string]Entry
}

func NewFileMapRepoMem() *RepoMem {
	return &RepoMem{
		fileMap: make(map[string]Entry),
	}
}

func (repo RepoMem) Find(path string) (Entry, error) {
	entry, ok := repo.fileMap[path]
	if ok {
		return entry, nil
	}
	return Entry{}, errors.New("not found")
}

func (repo RepoMem) Create(entry Entry) error {
	repo.fileMap[entry.Path] = entry
	return nil
}

func (repo RepoMem) Exists(path string) bool {
	_, ok := repo.fileMap[path]
	return ok
}

func (repo RepoMem) Delete(path string) error {
	delete(repo.fileMap, path)
	return nil
}

func (repo RepoMem) All() ([]Entry, error) {
	var entries = make([]Entry, 0)

	for _, entry := range repo.fileMap {
		entries = append(entries, entry)
	}

	return entries, nil
}

func (repo RepoMem) FindByID(id _id.ID) ([]Entry, error) {
	var entries = make([]Entry, 0)

	for _, entry := range repo.fileMap {
		if entry.ID == id {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}
