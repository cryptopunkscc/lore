package main

import (
	"gg/node"
	"log"
)

func main() {
	n, err := node.NewNode(node.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	err = n.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
