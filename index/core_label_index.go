package index

import (
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/story"
)

type CoreLabelIndex struct {
	repo CoreLabelRepo
}

func NewCoreLabelIndex(repo CoreLabelRepo) *CoreLabelIndex {
	return &CoreLabelIndex{repo: repo}
}

func (idx *CoreLabelIndex) Add(fileId string, source store.Reader) error {
	file, err := source.Read(fileId)
	if err != nil {
		return err
	}

	var coreLabel story.CoreLabel
	err = story.Parse(file, &coreLabel)
	if err != nil {
		return err
	}

	return idx.repo.Add(fileId, coreLabel)
}

func (idx *CoreLabelIndex) Remove(fileId string) error {
	return idx.repo.Remove(fileId)
}

func (idx *CoreLabelIndex) Query(q string) ([]string, error) {
	return idx.repo.Search(q)
}
