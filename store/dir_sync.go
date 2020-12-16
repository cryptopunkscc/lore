package store

import (
	"errors"
	"github.com/cryptopunkscc/lore/id"
	"os"
	"path/filepath"
)

// DirSync synchronizes the provided FileMapStore to a target directory
type DirSync struct {
	target       string
	fileMapStore *FileMapStore
	observer     Observer
}

func NewDirSync(target string, fileMapStore *FileMapStore, observer Observer) *DirSync {
	return &DirSync{
		target:       target,
		fileMapStore: fileMapStore,
		observer:     observer,
	}
}

// errAlreadyAdded means the path is already in the store
var errAlreadyAdded = errors.New("path already added")

// Update find and removes outdated entries in the FileMapStore and scans target directory for new files
func (sync *DirSync) Update() error {
	_ = sync.removeDirty()
	_ = sync.rescan()
	return nil
}

// removeDirty finds and removes outdated store entries
func (sync *DirSync) removeDirty() error {
	entries, err := sync.fileMapStore.Entries()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// We support only ID1 entries
		if !id.IsID1(entry.ID) {
			continue
		}

		// Get file size from the ID
		size, _ := id.ParseID1(entry.ID)

		// Get file info
		stat, err := os.Stat(entry.Path)
		if err != nil {
			// TODO: Decide: should we delete even temporarily inaccessible files?
			continue
		}

		// Validate size
		if uint64(stat.Size()) != size {
			sync.removeFile(entry)
			continue
		}

		// Validate modification time
		if stat.ModTime().After(entry.AddedAt) {
			sync.removeFile(entry)
			continue
		}
	}
	// All done!
	return nil
}

// removeFile removes a FileMapEntry from the store
func (sync *DirSync) removeFile(entry FileMapEntry) error {
	err := sync.fileMapStore.Remove(entry.Path)
	if err != nil {
		return err
	}

	// Notify the observer
	if sync.observer != nil {
		sync.observer.Removed(entry.ID)
	}

	return nil
}

// rescan scans target dir and adds any new files
func (sync *DirSync) rescan() error {
	err := filepath.Walk(sync.target, func(path string, info os.FileInfo, err error) error {
		// Skip erroneous entries
		if err != nil {
			return nil
		}

		// Skip empty files
		if info.Size() == 0 {
			return nil
		}

		if info.Mode().IsRegular() {
			_ = sync.addFile(path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// addFile adds a new file to the store
func (sync *DirSync) addFile(path string) error {
	// Check if the file is already tracked
	if sync.fileMapStore.Contains(path) {
		return errAlreadyAdded
	}

	// Resolve file's ID1
	fileId, err := id.ResolveFileID(path, id.NewID1Resolver())
	if err != nil {
		return err
	}

	// Map the file
	err = sync.fileMapStore.Add(path, fileId)
	if err != nil {
		return err
	}

	// Notify the observer
	if sync.observer != nil {
		sync.observer.Added(fileId)
	}

	return nil
}
