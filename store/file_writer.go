package store

import (
	"github.com/cryptopunkscc/lore/id"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Make sure *FileWriter satisfies Writer interface
var _ Writer = &FileWriter{}

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
func NewFileWriter(dir string, resolver id.Resolver) (*FileWriter, error) {
	// Create a temporary file
	tmpFile, err := ioutil.TempFile(dir, "tmp-")
	if err != nil {
		return nil, err
	}

	// Use default resolver if none provided
	if resolver == nil {
		resolver = id.DefaultResolver()
	}

	return &FileWriter{
		tmp:      tmpFile,
		resolver: resolver,
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
func (w *FileWriter) Finalize() (string, error) {
	// Close the temporary file
	tmpPath := w.tmp.Name()
	if err := w.tmp.Close(); err != nil {
		return "", err
	}

	// Resolve the resolver of the file
	fileId := w.resolver.Resolve()

	// Rename temporary file to its resolver
	dstPath := filepath.Join(w.dir, fileId)
	err := os.Rename(tmpPath, dstPath)
	if err != nil {
		return "", err
	}

	return fileId, nil
}

// Discard deletes the data
func (w *FileWriter) Discard() error {
	// Try to close, but ignore errors, since we want to delete the file anyways
	_ = w.tmp.Close()

	return os.Remove(w.tmp.Name())
}
