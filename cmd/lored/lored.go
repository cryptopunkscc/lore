package main

import (
	"github.com/cryptopunkscc/lore/node"
	"github.com/cryptopunkscc/lore/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

const defaultConfigPath = "~/.lore/config.yaml"

func main() {
	var config node.Config
	var configPath = defaultConfigPath

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	configPath, _ = util.ExpandPath(configPath)

	file, err := ioutil.ReadFile(configPath)
	if err == nil {
		err = yaml.Unmarshal(file, &config)
		if err != nil {
			log.Fatalln("Error reading config file:", err)
		}
	}

	n, err := node.NewNode(config)
	if err != nil {
		log.Fatalln(err)
	}

	err = n.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
