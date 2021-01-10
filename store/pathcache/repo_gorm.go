package pathcache

import (
	"fmt"
	_id "github.com/cryptopunkscc/lore/id"
	"gorm.io/gorm"
	"time"
)

// Check interface
var _ Repo = &RepoGorm{}

type RepoGorm struct {
	db *gorm.DB
}

// NewPathCacheRepoGorm returns a gorm implementation of Repo. It automatically creates necessary tables.
func NewPathCacheRepoGorm(db *gorm.DB) (*RepoGorm, error) {
	repo := &RepoGorm{
		db: db,
	}

	// Make sure db table is intact
	err := repo.db.AutoMigrate(&gormFileMapEntry{})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// Contains returns true if provided path exists in the repo, false otherwise
func (repo *RepoGorm) Exists(path string) bool {
	var count int64
	err := repo.db.Model(&gormFileMapEntry{}).Where("path = ?", path).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// Contains returns true if provided path exists in the repo, false otherwise
func (repo *RepoGorm) Find(path string) (Entry, error) {
	var entry gormFileMapEntry

	err := repo.db.Where("path = ?", path).First(&entry).Error
	if err != nil {
		return Entry{}, err
	}

	id, err := _id.Parse(entry.ID)
	if err != nil {
		return Entry{}, fmt.Errorf("error parsing database row: %w", err)
	}

	return Entry{
		Path:       entry.Path,
		ID:         id,
		VerifiedAt: entry.CreatedAt,
	}, nil
}

// Create adds a Entry to the database
func (repo *RepoGorm) Create(entry Entry) error {
	gormEntry := fileMapEntryToGormFileMapEntry(entry)
	return repo.db.Create(&gormEntry).Error
}

// Delete removes an entry from the database by path
func (repo *RepoGorm) Delete(path string) error {
	return repo.db.Delete(&gormFileMapEntry{Path: path}).Error
}

// All returns all Entry in the database
func (repo *RepoGorm) All() ([]Entry, error) {
	// Fetch the rows
	rows, err := repo.db.Find(&[]gormFileMapEntry{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Convert the results
	entries := make([]Entry, 0)
	for rows.Next() {
		var entry gormFileMapEntry
		_ = repo.db.ScanRows(rows, &entry)
		entries = append(entries, gormFileMapEntryToFileMapEntry(entry))
	}
	return entries, nil
}

// Paths returns a list of paths mapped to the provided ID
func (repo *RepoGorm) FindByID(id _id.ID) ([]Entry, error) {
	res := make([]Entry, 0)

	// Fetch all entries for the ID
	tx := repo.db.Where("id = ?", id.String()).Find(&[]gormFileMapEntry{})
	if tx.Error != nil {
		return nil, fmt.Errorf("database error: %w", tx.Error)
	}

	// Get rows
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}

	// Convert rows to []string
	for rows.Next() {
		var row gormFileMapEntry
		if err := repo.db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		i, err := _id.Parse(row.ID)
		if err != nil {
			return nil, fmt.Errorf("error parsing database row: %w", err)
		}
		res = append(res, Entry{
			Path:       row.Path,
			ID:         i,
			VerifiedAt: row.CreatedAt,
		})
	}

	return res, nil
}

type gormFileMapEntry struct {
	Path      string `gorm:"primaryKey"`
	ID        string `gorm:"index"`
	CreatedAt time.Time
}

func (gormFileMapEntry) TableName() string {
	return "file_map_entries"
}

func fileMapEntryToGormFileMapEntry(entry Entry) gormFileMapEntry {
	return gormFileMapEntry{
		ID:        entry.ID.String(),
		Path:      entry.Path,
		CreatedAt: entry.VerifiedAt,
	}
}

func gormFileMapEntryToFileMapEntry(entry gormFileMapEntry) Entry {
	i, _ := _id.Parse(entry.ID)
	return Entry{
		ID:         i,
		Path:       entry.Path,
		VerifiedAt: entry.CreatedAt,
	}
}
