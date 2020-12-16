package node

import (
	"github.com/cryptopunkscc/lore/util"
	"path/filepath"
)

const defaultNodeDir = "~/.lore"

type Config struct {
	NodeDir string
	Paths   []string
	Urls    []string
}

func (cfg Config) GetNodeDir() string {
	var res = defaultNodeDir

	if cfg.NodeDir != "" {
		res = cfg.NodeDir
	}

	res, _ = util.ExpandPath(res)
	return res
}

func (cfg Config) GetDataDir() string {
	return filepath.Join(cfg.GetNodeDir(), "data")
}
