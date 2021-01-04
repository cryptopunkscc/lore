package node

import (
	"fmt"
	"github.com/cryptopunkscc/lore/comm/server"
	"github.com/cryptopunkscc/lore/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	index       *NodeIndex
}

func (node *Node) Added(id string) {
	err := node.index.Add(id, node.deviceStore)
	if err != nil {
		fmt.Println(id, "indexing error:", err)
	}
}

func (node *Node) Removed(id string) {
	_ = node.index.Remove(id)
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
	node.db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set up device store
	node.deviceStore, err = NewDeviceStore(&node.config, node.db, node)
	if err != nil {
		return nil, err
	}

	// Set up the index
	node.index, err = NewNodeIndex(node.db)
	if err != nil {
		return nil, err
	}

	for _, path := range config.Paths {
		_ = node.deviceStore.AddLocalDir(path)
	}

	for _, url := range config.Urls {
		_ = node.deviceStore.AddNetworkStore(url)
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

	return node, nil
}

func (node *Node) Run() error {
	return node.server.Run()
}
