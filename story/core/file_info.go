package core

import (
	"github.com/cryptopunkscc/lore/story"
)

const FileInfoStoryType = "core.fileinfo"

type FileInfo struct {
	Story story.Header
	Name  string
	Type  string
}

func (info *FileInfo) Sanitize() {
	info.Story.Type = FileInfoStoryType
}

type FileInfoIndex interface {
	Add(id string, fileinfo FileInfo) error
	Remove(id string) error
	Search(query string) ([]string, error)
}
