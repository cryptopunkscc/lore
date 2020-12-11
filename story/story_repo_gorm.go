package story

import (
	"gorm.io/gorm"
)

var _ StoryRepo = &StoryRepoGorm{}

// StoryRepoGorm is a gorm implementation of StoryRepo
type StoryRepoGorm struct {
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

func NewStoryRepoGorm(db *gorm.DB) (*StoryRepoGorm, error) {
	var err error
	var repo = &StoryRepoGorm{db: db}

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

func (repo *StoryRepoGorm) SetStoryType(id string, typ string) error {
	t := &dbStoryType{
		ID:   id,
		Type: typ,
	}

	return repo.db.Create(&t).Error
}

func (repo *StoryRepoGorm) GetStoryType(id string) (string, error) {
	var t dbStoryType
	err := repo.db.Where("id = ?", id).First(&t).Error
	return t.ID, err
}

func (repo *StoryRepoGorm) SetStoryRels(id string, rels []string) error {
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

func (repo *StoryRepoGorm) GetStoryRels(id string) ([]string, error) {
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

func (repo *StoryRepoGorm) Forget(id string) error {
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
