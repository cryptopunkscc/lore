package server

import (
	"gg/lore/storage"
	"log"
)

type DownloadHandler struct {
	logger  *log.Logger
	storage *storage.LocalStorage
}

func (h *DownloadHandler) Handle(req *Request) {
	switch req.Method() {
	default:
		h.HandleDefault(req)
	}
}

func (h *DownloadHandler) HandleDefault(req *Request) {
	id := req.Method()

	// Check if we have the file
	path, err := h.storage.Path(id)
	if err != nil {
		req.NotFound()
		return
	}

	// Serve it!
	req.ServeFile(path)
}
