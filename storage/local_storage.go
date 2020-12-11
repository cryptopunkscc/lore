package storage

import (
	"fmt"
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/story"
	"github.com/cryptopunkscc/lore/story/core"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Make sure *LocalStorage satisfies Storage interface
var _ Storage = &LocalStorage{}

// LocalStorage represents a structure giving access to local files
type LocalStorage struct {
	rootDir      string
	locationRepo LocationRepo
	index        *Index
}

const defaultAppDir = ".lore"

// NewLocalStorage returns a new instance of LocalStorage. db is used for storing metadata about local files.
func NewLocalStorage(db *gorm.DB) (*LocalStorage, error) {
	s := &LocalStorage{
		locationRepo: newLocationDbRepo(db),
		rootDir:      defaultRootDir(),
	}

	s.index = NewIndex(db)

	_ = os.MkdirAll(s.dataDir(), 0700)
	return s, nil
}

// Create returns a writer that writes to local storage
func (s *LocalStorage) Create() (Writer, error) {
	return NewLocalStorageWriter(s.dataDir(), nil, s.locationRepo, s.index)
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

// AddStory adds a story to local storage
func (s *LocalStorage) AddStory(obj interface{}) error {
	// Marshal to YAML first
	bytes, err := yaml.Marshal(&obj)
	if err != nil {
		return err
	}

	// Create a file for the story in local storage
	w, err := s.Create()
	if err != nil {
		return err
	}

	// Write story to the file
	_, err = w.Write(bytes)
	if err != nil {
		return err
	}

	// Finalize and the final ID
	_, err = w.Finalize()
	if err != nil {
		return err
	}

	return nil
}

// Add adds info about an existing local file to the database
func (s *LocalStorage) Add(path string) (Location, error) {
	// Check if already added
	existing, err := s.locationRepo.Find(path)
	if err == nil {
		return existing, ErrAlreadyAdded
	}

	// ID the file
	fileId, err := id.ResolveFileID(path, id.NewID0Resolver())
	if err != nil {
		return Location{}, err
	}

	// Build Location object
	loc := Location{
		Location:   path,
		ID:         fileId,
		VerifiedAt: time.Now(),
	}

	// Save to database
	err = s.locationRepo.Create(loc)
	if err != nil {
		return Location{}, err
	}

	_, err = story.ParseHeaderFromFile(path)
	if err == nil {
		log.Println("added a story file")
		// Add to story index
		_ = s.index.AddFile(path)
	} else {
		// Add FileInfo
		info := core.FileInfo{
			Story: story.Header{
				Rel: []string{fileId},
			},
			Name: filepath.Base(path),
			Type: "",
		}

		mimeBytes, err := exec.Command("file", "--mime-type", "-L", path).Output()
		if err == nil {
			info.Type = strings.TrimSpace(string(mimeBytes[len(path)+2:]))
		}

		info.Sanitize()
		err = s.AddStory(info)
		if err != nil {
			log.Println("error generating FileInfo:", err)
		}
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

func (s *LocalStorage) Search(query string) map[string]string {
	var info core.FileInfo
	var res = make(map[string]string, 0)
	matches, _ := s.index.FileInfoIndex.Search(query)

	for _, m := range matches {
		r, _ := s.Open(m)
		data, _ := ioutil.ReadAll(r)
		_ = story.ParseStory(data, &info)
		for _, r := range info.Story.Rel {
			res[r] = info.Name
		}
	}

	return res
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
