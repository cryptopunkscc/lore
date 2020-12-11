package core

import (
	"gorm.io/gorm"
)

var _ FileInfoIndex = &FileInfoIndexGorm{}

type FileInfoIndexGorm struct {
	db *gorm.DB
}

type gormFileInfo struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"index"`
	Type string `gorm:"index"`
}

func (gormFileInfo) TableName() string { return "core_fileinfo_index" }

func NewFileInfoIndexGorm(db *gorm.DB) *FileInfoIndexGorm {
	idx := &FileInfoIndexGorm{db: db}

	_ = idx.db.AutoMigrate(&gormFileInfo{})

	return idx
}

func (idx *FileInfoIndexGorm) Add(id string, info FileInfo) error {
	var row gormFileInfo

	row.ID = id
	row.Name = info.Name
	row.Type = info.Type

	return idx.db.Create(&row).Error
}

func (idx *FileInfoIndexGorm) Remove(id string) error {
	return idx.db.Delete(&gormFileInfo{ID: id}).Error
}

func (idx *FileInfoIndexGorm) Search(query string) ([]string, error) {
	var row gormFileInfo
	var matches = make([]string, 0)

	rows, err := idx.db.Where("name LIKE ?", "%"+query+"%").Find(&[]gormFileInfo{}).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		_ = idx.db.ScanRows(rows, &row)
		matches = append(matches, row.ID)
	}

	return matches, nil
}
