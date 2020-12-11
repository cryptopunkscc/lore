package storage

import (
	"fmt"
	"github.com/cryptopunkscc/lore/id"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Make sure *LocalStorage satisfies Storage interface
var _ Storage = &LocalStorage{}

// LocalStorage represents a structure giving access to local files
type LocalStorage struct {
	rootDir      string
	locationRepo LocationRepo
}

const defaultAppDir = ".lore"

// NewLocalStorage returns a new instance of LocalStorage. db is used for storing metadata about local files.
func NewLocalStorage(db *gorm.DB) (*LocalStorage, error) {
	s := &LocalStorage{
		locationRepo: newLocationDbRepo(db),
		rootDir:      defaultRootDir(),
	}
	_ = os.MkdirAll(s.dataDir(), 0700)
	return s, nil
}

// Create returns a writer that writes to local storage
func (s *LocalStorage) Create() (Writer, error) {
	return NewLocalStorageWriter(s.dataDir(), nil, s.locationRepo)
}

// Delete all files with given ID from local storage
func (s *LocalStorage) Delete(id string) error {
	locations, err := s.locationRepo.FindByID(id)
	if err != nil {
		return err
	}
	for _, l := range locations {
		err = os.Remove(l.Location)
		if err != nil {
			return fmt.Errorf("can't delete file %s: %e", l.Location, err)
		}
		_ = s.locationRepo.Delete(l)
	}
	return nil
}

// Open opens a local file by resolver, returns ErrFileNotFound on error
func (s *LocalStorage) Open(id string) (ReadSeekCloser, error) {
	locations, err := s.locationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	for _, l := range locations {
		file, err := os.Open(l.Location)
		if err == nil {
			return file, nil
		}
	}
	return nil, ErrFileNotFound
}

// List returns the list of all unique IDs in local storage
func (s *LocalStorage) List() ([]string, error) {
	res := make([]string, 0)

	// Get all locations from the repository
	all, err := s.locationRepo.All()
	if err != nil {
		return res, err
	}

	// Remove duplicates
	seen := make(map[string]bool, 0)
	for _, l := range all {
		if seen[l.ID] != true {
			res = append(res, l.ID)
			seen[l.ID] = true
		}
	}

	return res, nil
}

// Contains checks whether an ID is available in the local storage
func (s *LocalStorage) Contains(id string) (bool, error) {
	locations, err := s.locationRepo.FindByID(id)
	if err != nil {
		return false, err
	}
	if len(locations) == 0 {
		return false, nil
	}
	return true, nil
}

// Add adds info about an existing local file to the database
func (s *LocalStorage) Add(path string) (Location, error) {
	// Get the absolute path first
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Location{}, err
	}

	// Check if already added
	existing, err := s.locationRepo.Find(absPath)
	if err == nil {
		return existing, ErrAlreadyAdded
	}

	// ID the file
	fileId, err := id.ResolveFileID(absPath, id.NewID0Resolver())
	if err != nil {
		return Location{}, err
	}

	// Build Location object
	loc := Location{
		Location:   absPath,
		ID:         fileId,
		VerifiedAt: time.Now(),
	}

	// Save to database
	err = s.locationRepo.Create(loc)
	if err != nil {
		return Location{}, err
	}

	return loc, nil
}

// Forget removes all info about a file from the database, but doesn't delete files
func (s *LocalStorage) Forget(id string) error {
	locations, err := s.locationRepo.FindByID(id)
	if err != nil {
		return err
	}
	for _, l := range locations {
		_ = s.locationRepo.Delete(l)
	}
	return nil
}

// Path returns a local path for provided id or ErrFileNotFound
func (s *LocalStorage) Path(id string) (string, error) {
	locations, err := s.locationRepo.FindByID(id)
	if err != nil {
		return "", err
	}

	fmt.Println("LOCS:", len(locations), locations)

	for _, l := range locations {
		_, err := os.Stat(l.Location)
		if err == nil {
			log.Println("Best location", l.Location)
			return l.Location, nil
		}
	}

	return "", ErrFileNotFound
}

// dataDir returns the data directory path
func (s *LocalStorage) dataDir() string {
	return filepath.Join(s.rootDir, "data")
}

// defaultRootDir returns the default directory for inventory
func defaultRootDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, defaultAppDir)
}
