package node

import (
	"fmt"
	"github.com/cryptopunkscc/lore/comm/server"
	"github.com/cryptopunkscc/lore/graph"
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
)

const dbFileName = "db.sqlite3"

var _ store.Observer = &Node{}

type Node struct {
	config      Config
	db          *gorm.DB
	deviceStore *DeviceStore
	server      *server.Server
	graph       *graph.Graph
}

func (node *Node) Added(id id.ID) {
	_, err := node.graph.AddFrom(id, node.deviceStore)
	if err != nil {
		fmt.Println(id, "error adding node to graph:", err)
	}
}

func (node *Node) Removed(id id.ID) {
	err := node.graph.Remove(id)
	if err != nil {
		fmt.Println(id, "error removing node from graph:", err)
	}
}

func NewNode(config Config) (*Node, error) {
	var err error

	node := &Node{
		config: config,
	}

	// Make sure node directory exists
	err = os.MkdirAll(node.config.GetNodeDir(), 0700)
	if err != nil {
		return nil, err
	}

	// Read the database
	dbPath := filepath.Join(node.config.GetNodeDir(), dbFileName)
	node.db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set up gorm graph repo
	graphRepo, err := graph.NewGraphRepoGorm(node.db)
	if err != nil {
		return nil, err
	}

	// Set up the graph
	node.graph, err = graph.NewGraph(graphRepo, nil)
	if err != nil {
		return nil, err
	}

	// Set up device store
	node.deviceStore, err = NewDeviceStore(&node.config, node.db, node)
	if err != nil {
		return nil, err
	}

	free, err := node.deviceStore.Free()
	if err != nil {
		log.Println("Error checking free space:", err)
	} else {
		log.Println("Free store space:", free)
	}

	// Instantiate the server
	node.server, err = server.NewServer(server.TCPContentConfig, node.deviceStore)
	if err != nil {
		return nil, err
	}

	node.deviceStore.Refresh()

	return node, nil
}

func (node *Node) Run() error {
	return node.server.Run()
}
