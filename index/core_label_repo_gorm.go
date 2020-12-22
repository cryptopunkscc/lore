package index

import (
	"github.com/cryptopunkscc/lore/story"
	"gorm.io/gorm"
)

const coreLabelTableName = "core_labels"

var _ CoreLabelRepo = &CoreLabelRepoGorm{}

type CoreLabelRepoGorm struct {
	db *gorm.DB
}

type gormCoreLabel struct {
	ID    string `gorm:"primaryKey"`
	Label string `gorm:"index"`
}

func (gormCoreLabel) TableName() string { return coreLabelTableName }

func NewCoreLabelRepoGorm(db *gorm.DB) (*CoreLabelRepoGorm, error) {
	idx := &CoreLabelRepoGorm{db: db}

	err := idx.db.AutoMigrate(&gormCoreLabel{})
	if err != nil {
		return nil, err
	}

	return idx, nil
}

func (idx *CoreLabelRepoGorm) Add(id string, coreLabel story.CoreLabel) error {
	var row gormCoreLabel

	row.ID = id
	row.Label = coreLabel.Label

	return idx.db.Create(&row).Error
}

func (idx *CoreLabelRepoGorm) Remove(id string) error {
	return idx.db.Delete(&gormCoreLabel{ID: id}).Error
}

func (idx *CoreLabelRepoGorm) Search(query string) ([]string, error) {
	var row gormCoreLabel
	var matches = make([]string, 0)

	rows, err := idx.db.Where("label LIKE ?", "%"+query+"%").Find(&[]gormCoreLabel{}).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		_ = idx.db.ScanRows(rows, &row)
		matches = append(matches, row.ID)
	}

	return matches, nil
}
