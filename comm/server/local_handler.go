package server

import (
	"github.com/cryptopunkscc/lore/lore/storage"
	"github.com/cryptopunkscc/lore/node/swarm"
	"log"
)

type LocalHandler struct {
	logger  *log.Logger
	storage *storage.LocalStorage
	swarm   *swarm.Swarm
}

func (h *LocalHandler) Handle(req *Request) {
	switch req.Method() {
	default:
		h.HandleDefault(req)
	}
}

func (h *LocalHandler) HandleDefault(req *Request) {
	id := req.Method()

	path, err := h.storage.Path(id)

	if err == nil {
		h.logger.Println("File", id, "found in local storage at", path)
		req.ServeFile(path)
		return
	}

	// Find source in my swarm
	client := h.swarm.FindSource(id)
	if client == nil {
		h.logger.Println("no source in swarm")
		req.NotFound()
		return
	}

	stream, err := client.Download().Stream(id)
	if err != nil {
		h.logger.Println("client error while streaming:", err)
		req.NotFound()
		return
	}

	req.Stream(stream)
}
