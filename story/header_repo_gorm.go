package story

import (
	"gorm.io/gorm"
)

var _ HeaderRepo = &HeaderRepoGorm{}

// HeaderRepoGorm is a gorm implementation of HeaderRepo
type HeaderRepoGorm struct {
	db *gorm.DB
}

type dbStoryType struct {
	ID   string `gorm:"primaryKey"`
	Type string `gorm:"index"`
}

func (dbStoryType) TableName() string { return "story_types" }

type dbStoryRel struct {
	ID  string `gorm:"primaryKey;index"`
	Rel string `gorm:"primaryKey;index"`
}

func (dbStoryRel) TableName() string { return "story_rels" }

func NewHeaderRepoGorm(db *gorm.DB) (*HeaderRepoGorm, error) {
	var err error
	var repo = &HeaderRepoGorm{db: db}

	err = repo.db.AutoMigrate(&dbStoryType{})
	if err != nil {
		return nil, err
	}

	err = repo.db.AutoMigrate(&dbStoryRel{})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *HeaderRepoGorm) Add(id string, header *Header) error {
	var err error

	err = repo.SetStoryType(id, header.Type)
	if err != nil {
		return err
	}

	return repo.SetStoryRels(id, header.Rel)
}

func (repo *HeaderRepoGorm) Remove(id string) error {
	var err error

	err = repo.db.Where("id = ?", id).Delete(&dbStoryRel{}).Error
	if err != nil {
		return err
	}

	err = repo.db.Where("id = ?", id).Delete(&dbStoryType{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *HeaderRepoGorm) SetStoryType(id string, typ string) error {
	t := &dbStoryType{
		ID:   id,
		Type: typ,
	}

	return repo.db.Create(&t).Error
}

func (repo *HeaderRepoGorm) GetStoryType(id string) (string, error) {
	var t dbStoryType
	err := repo.db.Where("id = ?", id).First(&t).Error
	return t.ID, err
}

func (repo *HeaderRepoGorm) SetStoryRels(id string, rels []string) error {
	for _, r := range rels {
		err := repo.db.Create(&dbStoryRel{
			ID:  id,
			Rel: r,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *HeaderRepoGorm) GetStoryRels(id string) ([]string, error) {
	var rels = make([]string, 0)

	rows, err := repo.db.Where("id = ?", id).Find(&[]dbStoryRel{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i dbStoryRel
		err = repo.db.ScanRows(rows, &i)
		if err != nil {
			return nil, err
		}
		rels = append(rels, i.Rel)
	}
	return rels, nil
}
