package server

import (
	"github.com/cryptopunkscc/lore/comm/proto"
	"github.com/cryptopunkscc/lore/lore/storage"
	"log"
)

type ItemHandler struct {
	logger  *log.Logger
	storage *storage.LocalStorage
}

func (h *ItemHandler) Handle(req *Request) {
	switch req.Method() {
	case "info":
		h.Info(req)
	default:
		req.NotFound()
	}
}

func (h *ItemHandler) Info(req *Request) {
	var err error
	var request proto.ItemInfoRequest
	var response proto.ItemInfoResponse

	// Parse the request
	_ = req.Unmarshal(&request)

	_, err = h.storage.Path(request.ID)
	if err != nil {
		req.NotFound()
		return
	}
	response.ID = request.ID
	_ = req.OK(&response)
}
