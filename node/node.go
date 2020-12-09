package node

import (
	"fmt"
	"gg/comm/server"
	"gg/lore/storage"
	"gg/node/swarm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

const dbPath = ".lore/db.sqlite3"

type Node struct {
	config       Config
	db           *gorm.DB
	localStorage *storage.LocalStorage
	swarm        *swarm.Swarm
	server       *server.Server
}

func NewNode(config Config) (*Node, error) {
	var err error

	node := &Node{
		config: config,
	}

	// Open the database
	home, _ := os.UserHomeDir()
	_ = os.MkdirAll(filepath.Join(home, ".lore"), 0700)
	absDbPath := filepath.Join(home, dbPath)
	node.db, err = gorm.Open(sqlite.Open(absDbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %e", err)
	}

	// Set up local storage
	node.localStorage, err = storage.NewLocalStorage(node.db)
	if err != nil {
		return nil, err
	}

	// Set up swarm
	node.swarm = swarm.NewSwarm()

	// Instantiate the server
	node.server, err = server.NewServer(server.TCPContentConfig, node.localStorage, node.swarm)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (node *Node) Run() error {
	return node.server.Run()
}
