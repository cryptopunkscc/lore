package storage

import (
	"errors"
	"os"
	"time"
)

type Location struct {
	Location   string
	ID         string
	VerifiedAt time.Time
}

type LocationRepo interface {
	Create(Location) error
	Update(Location) error
	Delete(Location) error
	All() ([]Location, error)
	FindByID(id string) ([]Location, error)
	Find(location string) (Location, error)
}

// ErrFileDirty file on the disk has been modified after its resolver was resolved
var ErrFileDirty = errors.New("file is dirty")

// Validate checks if location is still valid (file exists and wasn't modified)
func (loc *Location) Validate() error {
	stat, err := os.Stat(loc.Location)
	if err != nil {
		return err
	}
	if stat.ModTime().After(loc.VerifiedAt) {
		return ErrFileDirty
	}
	return nil
}
