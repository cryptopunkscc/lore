package server

import (
	"github.com/cryptopunkscc/lore/node/swarm"
	"github.com/cryptopunkscc/lore/storage"
	"log"
	"net"
	"net/http"
	"os"
)

type Server struct {
	cfg Config

	// services
	localStorage *storage.LocalStorage
	swarm        *swarm.Swarm
	logger       *log.Logger

	// handlers
	downloadHandler RequestHandler
	localHandler    RequestHandler
	itemHandler     RequestHandler
	adminHandler    RequestHandler

	// local
	lis       net.Listener
	mux       *http.ServeMux
	reqLogger *RequestLogger
}

type RequestHandler interface {
	Handle(request *Request)
}

const unixSocketPath = "/tmp/lore.sock"

func NewServer(cfg Config, localStorage *storage.LocalStorage, swarm *swarm.Swarm) (*Server, error) {
	srv := &Server{
		cfg:          cfg,
		localStorage: localStorage,
		swarm:        swarm,
		mux:          http.NewServeMux(),
	}

	// Set up logging
	srv.logger = log.New(os.Stderr, "", log.LstdFlags)
	srv.reqLogger = &RequestLogger{
		Logger: srv.logger,
		Prefix: "<http> ",
	}

	// Set up handlers
	srv.downloadHandler = &DownloadHandler{storage: srv.localStorage, logger: srv.logger}
	srv.localHandler = &LocalHandler{storage: srv.localStorage, logger: srv.logger, swarm: swarm}
	srv.itemHandler = &ItemHandler{storage: srv.localStorage, logger: srv.logger}
	srv.adminHandler = &AdminHandler{storage: srv.localStorage, swarm: swarm, logger: srv.logger}

	return srv, nil
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := Request{
		ResponseWriter: w,
		Request:        req,
	}

	// Log the raw HTTP request
	srv.reqLogger.ServeHTTP(w, req)

	// Handle the request
	switch r.Scope() {
	case "download":
		srv.downloadHandler.Handle(&r)
	case "local":
		srv.localHandler.Handle(&r)
	case "item":
		srv.itemHandler.Handle(&r)
	case "admin":
		srv.adminHandler.Handle(&r)
	default:
		r.NotFound()
	}
}

func (srv *Server) Run() error {
	var err error

	// Handle UNIX transport
	if srv.cfg.Transport == "unix" {
		log.Printf("Starting a new server on unix socket %s\n", unixSocketPath)
		srv.lis, err = net.Listen("unix", unixSocketPath)
		if err != nil {
			return err
		}
		defer srv.lis.Close()
		defer os.Remove(unixSocketPath)
	}

	// Handle TCP transport
	if srv.cfg.Transport == "tcp" {
		log.Printf("Starting a new server on tcp4 address %s\n", srv.cfg.Address)
		srv.lis, err = net.Listen("tcp4", srv.cfg.Address)
		if err != nil {
			return err
		}
		defer srv.lis.Close()
	}

	return http.Serve(srv.lis, srv)
}
