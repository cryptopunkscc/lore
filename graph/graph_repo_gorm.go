package graph

import (
	_id "github.com/cryptopunkscc/lore/id"
	"gorm.io/gorm"
)

var _ GraphRepo = &GraphRepoGorm{}

type GraphRepoGorm struct {
	db *gorm.DB
}

func NewGraphRepoGorm(db *gorm.DB) (*GraphRepoGorm, error) {
	var err error
	var repo = &GraphRepoGorm{db: db}

	err = repo.db.AutoMigrate(&gormNodeType{})
	if err != nil {
		return nil, err
	}

	err = repo.db.AutoMigrate(&gormEdge{})
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (repo *GraphRepoGorm) AddNode(node *Node) error {
	// TODO: Make this an atomic operation
	var err error

	// Save node type
	err = repo.db.Create(&gormNodeType{
		ID:      node.ID.String(),
		Type:    node.Type,
		SubType: node.SubType,
	}).Error
	if err != nil {
		return err
	}

	// Save edges
	for _, r := range node.Edges {
		err = repo.db.Create(&gormEdge{
			From: node.ID.String(),
			To:   r.String(),
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *GraphRepoGorm) RemoveNode(id _id.ID) error {
	var err error

	err = repo.db.Delete(&gormNodeType{ID: id.String()}).Error
	if err != nil {
		return err
	}

	err = repo.db.Where("'from' = ?", id).Delete(&gormEdge{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *GraphRepoGorm) FindNode(id _id.ID) (*Node, error) {
	var nodeType gormNodeType
	var node = &Node{}

	// Fetch types
	err := repo.db.Where("id = ?", id).First(&nodeType).Error
	if err != nil {
		return nil, err
	}

	node.Type = nodeType.Type
	node.SubType = nodeType.SubType

	if node.Type == TypeStory {
		node.Edges = make([]_id.ID, 0)

		rows, err := repo.db.Where("id = ?", id).Find(&[]gormEdge{}).Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var i gormEdge
			err = repo.db.ScanRows(rows, &i)
			if err != nil {
				return nil, err
			}

			toID, err := _id.Parse(i.To)
			if err != nil {
				return nil, err
			}

			node.Edges = append(node.Edges, toID)
		}
	}

	return node, err
}

func (repo *GraphRepoGorm) Objects(typ string) ([]string, error) {
	var list = make([]string, 0)

	q := repo.db.Where("type = ?", TypeObject)
	if typ != "" {
		q = q.Where("sub_type = ?", typ)
	}

	rows, err := q.Find(&[]gormNodeType{}).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var nodeType gormNodeType
		err = repo.db.ScanRows(rows, &nodeType)
		if err != nil {
			return nil, err
		}
		list = append(list, nodeType.ID)
	}

	return list, nil
}

func (repo *GraphRepoGorm) Stories(_ _id.ID, typ string) ([]string, error) {
	// TODO: Filter also by edge
	var list = make([]string, 0)

	q := repo.db.Where("type = ?", TypeStory)
	if typ != "" {
		q = q.Where("sub_type = ?", typ)
	}

	rows, err := q.Find(&[]gormNodeType{}).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var nodeType gormNodeType
		err = repo.db.ScanRows(rows, &nodeType)
		if err != nil {
			return nil, err
		}
		list = append(list, nodeType.ID)
	}

	return list, nil
}

type gormNodeType struct {
	ID      string `gorm:"primaryKey"`
	Type    string `gorm:"index"`
	SubType string `gorm:"index"`
}

func (gormNodeType) TableName() string {
	return "graph_node_types"
}

type gormEdge struct {
	From string `gorm:"primaryKey;index"`
	To   string `gorm:"primaryKey;index"`
}

func (gormEdge) TableName() string { return "graph_edges" }
