package index

import (
	"github.com/cryptopunkscc/lore/story"
	"gorm.io/gorm"
)

var _ StoryHeaderRepo = &StoryHeaderRepoGorm{}

// StoryHeaderRepoGorm is a gorm implementation of StoryHeaderRepo
type StoryHeaderRepoGorm struct {
	db *gorm.DB
}

func NewStoryHeaderRepoGorm(db *gorm.DB) (*StoryHeaderRepoGorm, error) {
	var err error
	var repo = &StoryHeaderRepoGorm{db: db}

	err = repo.db.AutoMigrate(&gormStoryHeaderType{})
	if err != nil {
		return nil, err
	}

	err = repo.db.AutoMigrate(&gormStoryHeaderRel{})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *StoryHeaderRepoGorm) Add(id string, header *story.Header) error {
	var err error

	err = repo.setStoryHeaderType(id, header.Type)
	if err != nil {
		return err
	}

	return repo.setStoryHeaderRels(id, header.Rel)
}

func (repo *StoryHeaderRepoGorm) Remove(id string) error {
	var err error

	err = repo.db.Where("id = ?", id).Delete(&gormStoryHeaderRel{}).Error
	if err != nil {
		return err
	}

	err = repo.db.Where("id = ?", id).Delete(&gormStoryHeaderType{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *StoryHeaderRepoGorm) Find(types []string, rels []string) ([]string, error) {
	q := repo.db

	if types != nil {
		q = q.Where("type in ?", types)
	}

	return nil, nil
}

func (repo *StoryHeaderRepoGorm) setStoryHeaderType(id string, typ string) error {
	t := &gormStoryHeaderType{
		ID:   id,
		Type: typ,
	}

	return repo.db.Create(&t).Error
}

func (repo *StoryHeaderRepoGorm) getStoryHeaderType(id string) (string, error) {
	var t gormStoryHeaderType
	err := repo.db.Where("id = ?", id).First(&t).Error
	return t.ID, err
}

func (repo *StoryHeaderRepoGorm) setStoryHeaderRels(id string, rels []string) error {
	for _, r := range rels {
		err := repo.db.Create(&gormStoryHeaderRel{
			ID:  id,
			Rel: r,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *StoryHeaderRepoGorm) getStoryHeaderRels(id string) ([]string, error) {
	var rels = make([]string, 0)

	rows, err := repo.db.Where("id = ?", id).Find(&[]gormStoryHeaderRel{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i gormStoryHeaderRel
		err = repo.db.ScanRows(rows, &i)
		if err != nil {
			return nil, err
		}
		rels = append(rels, i.Rel)
	}
	return rels, nil
}

type gormStoryHeaderType struct {
	ID   string `gorm:"primaryKey"`
	Type string `gorm:"index"`
}

func (gormStoryHeaderType) TableName() string { return "story_header_types" }

type gormStoryHeaderRel struct {
	ID  string `gorm:"primaryKey;index"`
	Rel string `gorm:"primaryKey;index"`
}

func (gormStoryHeaderRel) TableName() string { return "story_header_rels" }
