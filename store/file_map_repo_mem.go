package store

import "time"

var _ FileMapRepo = &FileMapRepoMem{}

type FileMapRepoMem struct {
	fileMap map[string]memFileMapEntry
}

func NewFileMapRepoMem() *FileMapRepoMem {
	return &FileMapRepoMem{
		fileMap: make(map[string]memFileMapEntry),
	}
}

func (repo FileMapRepoMem) Contains(path string) bool {
	_, ok := repo.fileMap[path]

	return ok
}

func (repo FileMapRepoMem) Create(entry FileMapEntry) error {
	repo.fileMap[entry.Path] = memFileMapEntry{
		ID:      entry.ID,
		AddedAt: time.Now(),
	}
	return nil
}

func (repo FileMapRepoMem) Delete(path string) error {
	delete(repo.fileMap, path)
	return nil
}

func (repo FileMapRepoMem) All() ([]FileMapEntry, error) {
	var entries = make([]FileMapEntry, 0)

	for k, v := range repo.fileMap {
		entries = append(entries, FileMapEntry{
			Path:    k,
			ID:      v.ID,
			AddedAt: v.AddedAt,
		})
	}

	return entries, nil
}

func (repo FileMapRepoMem) Paths(id string) ([]string, error) {
	var paths = make([]string, 0)

	for path, info := range repo.fileMap {
		if info.ID == id {
			paths = append(paths, path)
		}
	}

	return paths, nil
}

type memFileMapEntry struct {
	ID      string
	AddedAt time.Time
}
