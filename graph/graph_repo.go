package graph

const (
	TypeObject = "object"
	TypeStory  = "story"
)

type GraphRepo interface {
	AddNode(node *Node) error
	RemoveNode(id string) error
	FindNode(id string) (*Node, error)
	Stories(edge string, typ string) ([]string, error)
	Objects(typ string) ([]string, error)
}
