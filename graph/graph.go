package graph

import (
	"bytes"
	"fmt"
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/story"
	"github.com/cryptopunkscc/lore/util"
)

type Graph struct {
	repo GraphRepo
}

type Node struct {
	ID      string
	Type    string
	SubType string
	Edges   []string
}

func NewGraph(repo GraphRepo) (*Graph, error) {
	return &Graph{repo: repo}, nil
}

// Add a node to the graph
func (graph *Graph) Add(nodeId string, reader store.Reader) (*Node, error) {
	node := &Node{ID: nodeId}

	file, err := reader.Read(nodeId)
	if err != nil {
		return nil, err
	}

	var buf = make([]byte, story.MaxStorySize)
	n, err := file.Read(buf)
	if err != nil {
		return nil, err
	}

	header, err := story.ParseHeaderFromBytes(buf[:n])
	if err != nil {
		fmt.Println(err)
		node.Type = TypeObject
		node.SubType, err = util.GetContentType(bytes.NewReader(buf))
		if err != nil {
			return nil, err
		}
	} else {
		node.Type = TypeStory
		node.SubType = header.Type
		node.Edges = header.Rel
	}

	err = graph.repo.AddNode(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// Remove a node from the graph
func (graph *Graph) Remove(id string) error {
	return graph.repo.RemoveNode(id)
}

// Get returns info about a node in the graph
func (graph *Graph) Get(id string) (*Node, error) {
	return graph.repo.FindNode(id)
}

// Objects returns a list of all objects in the graph. If typ is non-zero, only objects of the given type will be returned.
func (graph *Graph) Objects(typ string) ([]string, error) {
	return graph.repo.Objects(typ)
}

// Stories returns a list of stories. If rel is non-zero, the list will be limited to stories that are related to the provided id.
// If typ is non-zero, results will be limited to the provided type.
func (graph *Graph) Stories(edge string, typ string) ([]string, error) {
	return graph.repo.Stories(edge, typ)
}
