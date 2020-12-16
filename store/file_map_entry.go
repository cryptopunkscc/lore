package store

import (
	"time"
)

type FileMapEntry struct {
	Path    string
	ID      string
	AddedAt time.Time
}
