package storage

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/storage/index"
	"path/filepath"
	"time"
)

type LocalStorageWriter struct {
	writer       Writer
	dir          string
	locationRepo LocationRepo
	storyIndex   *index.StoryIndex
}

func NewLocalStorageWriter(dir string, resolver id.Resolver, locRepo LocationRepo, storyIndex *index.StoryIndex) (*LocalStorageWriter, error) {
	fw, err := NewFileWriter(dir, resolver)
	if err != nil {
		return nil, err
	}

	w := &LocalStorageWriter{
		writer:       fw,
		dir:          dir,
		locationRepo: locRepo,
		storyIndex:   storyIndex,
	}
	return w, nil
}

func (w *LocalStorageWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w *LocalStorageWriter) Discard() error {
	return w.writer.Discard()
}

func (w *LocalStorageWriter) Finalize() (string, error) {
	// Finalize writing data
	fileId, err := w.writer.Finalize()
	if err != nil {
		return fileId, err
	}

	// Get the path to newly written file
	path, err := filepath.Abs(filepath.Join(w.dir, fileId))
	if err != nil {
		return "", err
	}

	// Write file info to the repository
	l := Location{
		Location:   path,
		ID:         fileId,
		VerifiedAt: time.Now(),
	}
	err = w.locationRepo.Create(l)
	if err != nil {
		return "", err
	}

	_ = w.storyIndex.IndexFile(path)

	return fileId, nil
}
