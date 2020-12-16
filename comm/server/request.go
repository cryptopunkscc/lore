package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func (req *Request) Unmarshal(target interface{}) error {
	var err error

	// Read the bytes
	bodyBytes, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, target)
}

func (req *Request) OK(data interface{}) error {
	var err error

	// Marshal the response to a JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the response
	_, err = req.ResponseWriter.Write(jsonBytes)
	return err
}

func (req *Request) ServerError(msg string) error {
	var err error
	req.ResponseWriter.WriteHeader(http.StatusInternalServerError)
	_, err = req.ResponseWriter.Write([]byte(msg))
	return err
}

func (req *Request) NotFound() {
	http.NotFound(req.ResponseWriter, req.Request)
}

func (req *Request) ServeFile(path string) {
	http.ServeFile(req.ResponseWriter, req.Request, path)
}

func (req *Request) Serve(name string, content io.ReadSeeker) {
	http.ServeContent(req.ResponseWriter, req.Request, name, time.Time{}, content)
}

func (req *Request) Scope() string {
	segments := strings.Split(req.Request.URL.Path, "/")[1:]
	return segments[0]
}

func (req *Request) Method() string {
	segments := strings.Split(req.Request.URL.Path, "/")[1:]
	if len(segments) < 2 {
		return ""
	}
	return segments[1]
}

func (req *Request) Stream(reader io.Reader) {
	req.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	_, _ = io.Copy(req.ResponseWriter, reader)
}
