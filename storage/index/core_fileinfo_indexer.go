package index

import (
	"github.com/cryptopunkscc/lore/story"
	"gorm.io/gorm"
)

var _ Indexer = &coreFileinfoIndexer{}

type coreFileinfo struct {
	Name string
	Type string
}

type coreFileinfoIndexer struct {
	db *gorm.DB
}

func (idx *coreFileinfoIndexer) Query(query string) ([]string, error) {
	var row dbCoreFileinfo
	var matches = make([]string, 0)

	rows, err := idx.db.Where("name LIKE ?", "%"+query+"%").Find(&[]dbCoreFileinfo{}).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		_ = idx.db.ScanRows(rows, &row)
		matches = append(matches, row.ID)
	}

	return matches, nil
}

func (idx *coreFileinfoIndexer) Deindex(id string) error {
	return idx.db.Delete(&dbCoreFileinfo{ID: id}).Error
}

type dbCoreFileinfo struct {
	ID   string `gorm:"primaryKey"`
	Name string `gorm:"index"`
	Type string `gorm"index"`
}

func (dbCoreFileinfo) TableName() string { return "core_fileinfo_index" }

func newCoreFileinfoIndexer(db *gorm.DB) *coreFileinfoIndexer {
	idx := &coreFileinfoIndexer{db: db}

	_ = idx.db.AutoMigrate(&dbCoreFileinfo{})

	return idx
}

func (idx *coreFileinfoIndexer) Index(id string, header *story.Header, data []byte) error {
	var err error
	var s coreFileinfo
	var row dbCoreFileinfo

	err = story.ParseStory(data, &s)
	if err != nil {
		return err
	}

	row.ID = id
	row.Name = s.Name
	row.Type = s.Type

	idx.db.Create(&row)

	return nil
}
