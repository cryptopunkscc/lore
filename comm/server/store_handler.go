package server

import (
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
	id := req.Method()

	// Open the file
	file, err := handler.store.Read(id)
	if err != nil {
		req.ServerError(err.Error())
		return
	}

	// Serve it!
	req.Serve(id, file)
}

func (handler *StoreHandler) HandleList(req *Request) {
	list, err := handler.store.List()
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

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

	_, _ = req.ResponseWriter.Write([]byte(id))
}

func (handler *StoreHandler) HandleDelete(req *Request) {
	id := req.Method()

	err := handler.store.Delete(id)
	if err != nil {
		_ = req.ServerError(err.Error())
		return
	}

	req.OK(nil)
}

func (handler *StoreHandler) HandleFree(req *Request) {
	free, _ := handler.store.Free()

	req.OK(free)
}
