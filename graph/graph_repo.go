package graph

import "github.com/cryptopunkscc/lore/id"

const (
	TypeObject = "object"
	TypeStory  = "story"
)

type GraphRepo interface {
	AddNode(node *Node) error
	RemoveNode(id id.ID) error
	FindNode(id id.ID) (*Node, error)
	Stories(node id.ID, typ string) ([]string, error)
	Objects(typ string) ([]string, error)
}
