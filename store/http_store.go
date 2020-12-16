package store

import (
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/lore/httpfile"
	"io"
	"io/ioutil"
	"net/http"
	_url "net/url"
	"path/filepath"
)

var _ ReadEditor = &HTTPStore{}

type HTTPStore struct {
	baseUrl string
}

func NewHTTPStore(baseUrl string) *HTTPStore {
	return &HTTPStore{baseUrl: baseUrl}
}

func (s HTTPStore) Read(id string) (ReadSeekCloser, error) {
	url := s.url(id)

	return httpfile.Open(url)
}

func (s HTTPStore) List() ([]string, error) {
	res, err := http.Get(s.url(""))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var list = make([]string, 0)
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s HTTPStore) Create() (Writer, error) {
	return newHttpWriter(s.url(""))
}

func (s HTTPStore) Delete(id string) error {
	req, err := http.NewRequest("DELETE", s.url(id), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http status %d", res.StatusCode)
	}

	return nil
}

func (s HTTPStore) url(path string) string {
	url, _ := _url.Parse(s.baseUrl)
	url.Path = filepath.Join(url.Path, path)
	return url.String()
}

var _ Writer = &httpWriter{}

type httpWriter struct {
	req *http.Request
	res chan httpFinalizeResponse
	r   *io.PipeReader
	w   *io.PipeWriter
}

type httpFinalizeResponse struct {
	id  string
	err error
}

func newHttpWriter(url string) (*httpWriter, error) {
	var err error
	var w = &httpWriter{
		res: make(chan httpFinalizeResponse),
	}

	w.req, err = http.NewRequest("POST", url, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (writer *httpWriter) Finalize() (string, error) {
	_ = writer.w.Close()
	res := <-writer.res
	return res.id, res.err
}

func (writer *httpWriter) Discard() error {
	return ErrUnsupported
}

func (writer *httpWriter) Write(data []byte) (int, error) {
	if writer.w == nil {
		writer.r, writer.w = io.Pipe()
		go writer.execute()
	}

	return writer.w.Write(data)
}

func (writer *httpWriter) Read(data []byte) (int, error) {
	return writer.r.Read(data)
}

func (writer *httpWriter) execute() {
	res, err := http.DefaultClient.Do(writer.req)
	if err != nil {
		writer.res <- httpFinalizeResponse{"", err}
		return
	}

	id, err := ioutil.ReadAll(res.Body)
	if err != nil {
		writer.res <- httpFinalizeResponse{"", err}
		return
	}

	writer.res <- httpFinalizeResponse{string(id), nil}
	return
}
