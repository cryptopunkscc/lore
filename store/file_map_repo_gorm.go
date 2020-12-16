package store

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Check interface
var _ FileMapRepo = &FileMapRepoGorm{}

type FileMapRepoGorm struct {
	db *gorm.DB
}

// NewFileMapRepoGorm returns a gorm implementation of FileMapRepo. It automatically creates necessary tables.
func NewFileMapRepoGorm(db *gorm.DB) (*FileMapRepoGorm, error) {
	repo := &FileMapRepoGorm{
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
func (repo *FileMapRepoGorm) Contains(path string) bool {
	var count int64
	err := repo.db.Model(&gormFileMapEntry{}).Where("path = ?", path).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// Create adds a FileMapEntry to the database
func (repo *FileMapRepoGorm) Create(entry FileMapEntry) error {
	gormEntry := fileMapEntryToGormFileMapEntry(entry)
	return repo.db.Create(&gormEntry).Error
}

// Delete removes an entry from the database by path
func (repo *FileMapRepoGorm) Delete(path string) error {
	return repo.db.Delete(&gormFileMapEntry{Path: path}).Error
}

// All returns all FileMapEntry in the database
func (repo *FileMapRepoGorm) All() ([]FileMapEntry, error) {
	// Fetch the rows
	rows, err := repo.db.Find(&[]gormFileMapEntry{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Convert the results
	entries := make([]FileMapEntry, 0)
	for rows.Next() {
		var entry gormFileMapEntry
		repo.db.ScanRows(rows, &entry)
		entries = append(entries, gormFileMapEntryToFileMapEntry(entry))
	}
	return entries, nil
}

// Paths returns a list of paths mapped to the provided ID
func (repo *FileMapRepoGorm) Paths(id string) ([]string, error) {
	res := make([]string, 0)

	// Fetch all entries for the ID
	tx := repo.db.Where("id = ?", id).Find(&[]gormFileMapEntry{})
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
		var entry gormFileMapEntry
		if err := repo.db.ScanRows(rows, &entry); err != nil {
			return nil, err
		}
		res = append(res, entry.Path)
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

func fileMapEntryToGormFileMapEntry(i FileMapEntry) gormFileMapEntry {
	return gormFileMapEntry{
		ID:        i.ID,
		Path:      i.Path,
		CreatedAt: i.AddedAt,
	}
}

func gormFileMapEntryToFileMapEntry(i gormFileMapEntry) FileMapEntry {
	return FileMapEntry{
		ID:      i.ID,
		Path:    i.Path,
		AddedAt: i.CreatedAt,
	}
}
