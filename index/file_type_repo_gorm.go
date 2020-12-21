package index

import "gorm.io/gorm"

var _ FileTypeRepo = &FileTypeRepoGorm{}

type FileTypeRepoGorm struct {
	db *gorm.DB
}

type gormFileType struct {
	ID   string `gorm:"primaryKey"`
	Type string `gorm:"index"`
}

func (gormFileType) TableName() string {
	return "file_types"
}

func NewFileTypeRepoGorm(db *gorm.DB) (*FileTypeRepoGorm, error) {
	repo := &FileTypeRepoGorm{db: db}

	err := repo.db.AutoMigrate(&gormFileType{})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo FileTypeRepoGorm) Get(id string) string {
	var res gormFileType
	err := repo.db.Where("id = ?", id).First(&res).Error
	if err != nil {
		return ""
	}
	return res.Type
}

func (repo *FileTypeRepoGorm) Set(id string, typ string) error {
	return repo.db.Create(&gormFileType{
		ID:   id,
		Type: typ,
	}).Error
}

func (repo *FileTypeRepoGorm) Clear(id string) error {
	return repo.db.Delete(&gormFileType{ID: id}).Error
}
