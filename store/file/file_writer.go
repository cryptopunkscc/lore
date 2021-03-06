package file

import (
	"errors"
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Make sure *FileWriter satisfies Writer interface
var _ store.Writer = &FileWriter{}

// FileWriter writes data to a local file and once all data is written, it resolves the ID of the file,
// renames the file to that ID and returns the ID.
type FileWriter struct {
	tmp      *os.File
	resolver id.Resolver
	dir      string
}

// NewFileWriter returns a writer that writes to a temporary file in the provided directory.
// Call Discard() to stop writing and delete the data written so far.
// Call Finalize() to close the file and resolve its final ID. The file will remain in the directory, but will be
//   renamed to its ID.
func NewFileWriter(dir string) (*FileWriter, error) {
	// Create a temporary file
	tmpFile, err := ioutil.TempFile(dir, "tmp-")
	if err != nil {
		return nil, err
	}

	return &FileWriter{
		tmp:      tmpFile,
		resolver: id.NewResolver(),
		dir:      dir,
	}, nil
}

// Write writes more data to the file
func (w *FileWriter) Write(data []byte) (int, error) {
	i, err := w.resolver.Write(data)
	if err != nil {
		return i, err
	}
	return w.tmp.Write(data)
}

// Finalize closes the file, renames it to its resolver and returns the resolver
func (w *FileWriter) Finalize() (id.ID, error) {
	var err error

	// Close the temporary file
	tmpPath := w.tmp.Name()
	if err := w.tmp.Close(); err != nil {
		return id.ID{}, err
	}

	// Resolve the resolver of the file
	fileId := w.resolver.Resolve()

	dstPath := filepath.Join(w.dir, fileId.String())

	_, err = os.Stat(dstPath)
	if !errors.Is(err, os.ErrNotExist) {
		_ = w.Discard()
		return fileId, store.ErrAlreadyExists
	}

	// Rename temporary file to its id
	err = os.Rename(tmpPath, dstPath)
	if err != nil {
		return id.ID{}, err
	}

	return fileId, nil
}

// Discard deletes the data
func (w *FileWriter) Discard() error {
	// Try to close, but ignore errors, since we want to delete the file anyways
	_ = w.tmp.Close()

	return os.Remove(w.tmp.Name())
}
