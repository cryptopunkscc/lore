package node

import (
	"fmt"
	"github.com/cryptopunkscc/lore/comm/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

const dbFileName = "db.sqlite3"

type Node struct {
	config      Config
	db          *gorm.DB
	deviceStore *DeviceStore
	server      *server.Server
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
	node.deviceStore, err = NewDeviceStore(&node.config, node.db)
	if err != nil {
		return nil, err
	}

	for _, path := range config.Paths {
		_ = node.deviceStore.AddLocalDir(path)
	}

	for _, url := range config.Urls {
		_ = node.deviceStore.AddNetworkStore(url)
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
