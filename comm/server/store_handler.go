package server

import (
	_id "github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"io"
	"log"
	"net/http"
)

type StoreHandler struct {
	logger *log.Logger
	store  store.Store
}

func (handler *StoreHandler) Handle(req *Request) {
	switch req.Method() {
	case "":
		switch req.Request.Method {
		case http.MethodGet:
			handler.HandleList(req)
		case http.MethodPost:
			handler.HandleCreate(req)
		}
	case "free":
		if req.Request.Method == http.MethodGet {
			handler.HandleFree(req)
		}
	default:
		switch req.Request.Method {
		case http.MethodGet:
			handler.HandleRead(req)
		case http.MethodDelete:
			handler.HandleDelete(req)
		}
	}
}

func (handler *StoreHandler) HandleRead(req *Request) {
	idStr := req.Method()

	id, err := _id.Parse(idStr)
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	// Open the file
	file, err := handler.store.Read(id)
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	// Serve it!
	req.Serve(id.String(), file)
}

func (handler *StoreHandler) HandleList(req *Request) {
	set, err := handler.store.List()
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	// Convert to string list
	var list = make([]string, 0)
	set.Each(func(id _id.ID) {
		list = append(list, id.String())
	})

	_ = req.OK(list)
}

func (handler *StoreHandler) HandleCreate(req *Request) {
	writer, err := handler.store.Create()
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	_, err = io.Copy(writer, req.Request.Body)
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	id, err := writer.Finalize()
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	_, _ = req.ResponseWriter.Write([]byte(id.String()))
}

func (handler *StoreHandler) HandleDelete(req *Request) {
	idStr := req.Method()

	id, err := _id.Parse(idStr)
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	err = handler.store.Delete(id)
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	_ = req.OK(nil)
}

func (handler *StoreHandler) HandleFree(req *Request) {
	free, _ := handler.store.Free()

	_ = req.OK(free)
}
