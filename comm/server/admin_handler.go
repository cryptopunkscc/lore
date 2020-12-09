package server

import (
	"github.com/cryptopunkscc/lore/comm/proto"
	"github.com/cryptopunkscc/lore/lore/storage"
	"github.com/cryptopunkscc/lore/node/swarm"
	"log"
)

type AdminHandler struct {
	logger  *log.Logger
	storage *storage.LocalStorage
	swarm   *swarm.Swarm
}

func (admin *AdminHandler) Handle(req *Request) {
	switch req.Method() {
	case "add":
		admin.Add(req)
	case "list":
		admin.List(req)
	case "listsources":
		admin.ListSources(req)
	case "addsource":
		admin.AddSource(req)
	case "removesource":
		admin.RemoveSource(req)
	default:
		req.NotFound()
	}
}

func (admin *AdminHandler) Add(req *Request) {
	var err error
	var addRequest proto.AdminAddRequest
	var addResponse proto.AdminAddResponse

	err = req.Unmarshal(&addRequest)
	if err != nil {
		return
	}

	item, err := admin.storage.Add(addRequest.Path)

	if err != nil {
		addResponse.Error = err.Error()
	} else {
		addResponse.ID = item.ID
		addResponse.Size = int(0)
		addResponse.Type = ""
		addResponse.SubType = ""
	}

	_ = req.OK(addResponse)
}

func (admin *AdminHandler) List(req *Request) {
	var err error
	var listResponse proto.AdminListResponse

	list, err := admin.storage.List()
	if err != nil {
		listResponse.Error = err.Error()
		_ = req.OK(&listResponse)
	}

	listResponse.Items = list
	_ = req.OK(&listResponse)
}

func (admin *AdminHandler) AddSource(req *Request) {
	var err error
	var request proto.AdminAddSourceRequest
	var response proto.AdminAddSourceResponse

	err = req.Unmarshal(&request)
	if err != nil {
		req.NotFound()
		return
	}

	err = admin.swarm.Add(request.Address)
	if err != nil {
		response.Error = err.Error()
	}

	_ = req.OK(&response)
}

func (admin *AdminHandler) RemoveSource(req *Request) {
	var err error
	var request proto.AdminRemoveSourceRequest
	var response proto.AdminRemoveSourceResponse

	err = req.Unmarshal(&request)
	if err != nil {
		req.NotFound()
		return
	}

	admin.swarm.Remove(request.Address)

	_ = req.OK(&response)
}

func (admin *AdminHandler) ListSources(req *Request) {
	var response proto.AdminListSourcesResponse

	response.Sources = admin.swarm.List()

	_ = req.OK(&response)
}
