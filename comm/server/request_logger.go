package server

import (
	"log"
	"net/http"
)

// Make sure *RequestLogger satisfies http.Handler interface
var _ http.Handler = &RequestLogger{}

// RequestLogger logs request info to provided Logger and passes the request on to the handler
type RequestLogger struct {
	Handler http.Handler
	Logger  *log.Logger
	Prefix  string
}

// ServeHTTP handles an incoming HTTP request. It will log request info if Logger is present, and then
// pass the request to the next http.Handler if present.
func (l *RequestLogger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if l.Logger != nil {
		l.Logger.Printf("%s[%s] %s %s\n", l.Prefix, req.RemoteAddr, req.Method, req.URL.Path)
	}
	if l.Handler != nil {
		l.Handler.ServeHTTP(w, req)
	}
}
