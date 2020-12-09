package storage

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

// Check interface
var _ = (LocationRepo)(&locationDbRepo{})

type dbLocation struct {
	gorm.Model
	Location   string `gorm:"primaryKey"`
	ID         string `gorm:"index"`
	VerifiedAt time.Time
}

func (dbLocation) TableName() string {
	return "locations"
}

type locationDbRepo struct {
	db *gorm.DB
}

func newLocationDbRepo(db *gorm.DB) *locationDbRepo {
	repo := &locationDbRepo{
		db: db,
	}
	_ = repo.db.AutoMigrate(&dbLocation{})
	return repo
}

func (repo *locationDbRepo) Create(loc Location) error {
	_loc := locationToDbLocation(loc)
	return repo.db.Create(&_loc).Error
}

func (repo *locationDbRepo) Update(loc Location) error {
	_loc := locationToDbLocation(loc)
	return repo.db.Save(&_loc).Error
}

func (repo *locationDbRepo) Delete(loc Location) error {
	_loc := locationToDbLocation(loc)
	return repo.db.Delete(&_loc).Error
}

func (repo *locationDbRepo) CreateOrUpdate(loc Location) error {
	_loc := locationToDbLocation(loc)
	if repo.Exists(loc) {
		return repo.db.Save(&_loc).Error
	}
	return repo.db.Create(&_loc).Error
}

func (repo *locationDbRepo) All() ([]Location, error) {
	rows, err := repo.db.Find(&[]dbLocation{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]Location, 0)

	for rows.Next() {
		var i dbLocation
		repo.db.ScanRows(rows, &i)
		files = append(files, dbLocationToLocation(i))
	}
	return files, nil
}

func (repo *locationDbRepo) Exists(loc Location) bool {
	return repo.db.Where("location = ?", loc.Location).First(&dbLocation{}).RowsAffected > 0
}

func (repo *locationDbRepo) FindByID(id string) ([]Location, error) {
	res := make([]Location, 0)

	var arr []dbLocation

	// Fetch all locations for the ID
	tx := repo.db.Where("id = ?", id).Find(&arr)
	if tx.Error != nil {
		return nil, fmt.Errorf("database error: %e", tx.Error)
	}

	// Get rows
	rows, err := tx.Rows()
	if err != nil {
		return nil, err
	}

	// Convert rows to model
	for rows.Next() {
		var loc dbLocation
		if err := repo.db.ScanRows(rows, &loc); err != nil {
			return nil, err
		}
		res = append(res, dbLocationToLocation(loc))
	}

	return res, nil
}

func (repo *locationDbRepo) Find(path string) (Location, error) {
	var loc dbLocation
	tx := repo.db.Where("location = ?", path).First(&loc)
	return dbLocationToLocation(loc), tx.Error
}

// conversions
func locationToDbLocation(i Location) dbLocation {
	return dbLocation{
		ID:         i.ID,
		Location:   i.Location,
		VerifiedAt: i.VerifiedAt,
	}
}

func dbLocationToLocation(i dbLocation) Location {
	return Location{
		ID:         i.ID,
		Location:   i.Location,
		VerifiedAt: i.VerifiedAt,
	}
}
