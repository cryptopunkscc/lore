package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cryptopunkscc/lore/httpfile"
	"github.com/tv42/httpunix"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type HTTPClient struct {
	logger *log.Logger
	url    string
	http   http.Client
}

func NewHTTPClient(rootURL string) *HTTPClient {
	u, err := url.Parse(rootURL)
	if err != nil {
		return nil
	}

	var client = http.Client{}

	if u.Scheme == "http+unix" {
		t := &httpunix.Transport{
			DialTimeout:           100 * time.Millisecond,
			RequestTimeout:        1 * time.Second,
			ResponseHeaderTimeout: 1 * time.Second,
		}
		t.RegisterLocation("lored", "/tmp/lore.sock")
		client.Transport = t
	}

	return &HTTPClient{
		url:  rootURL,
		http: client,
	}
}

func (client *HTTPClient) GetStream(scope string, method string) (io.ReadSeeker, error) {
	// Prepare full URL
	u, _ := url.Parse(client.url)
	u.Path = filepath.Join(u.Path, scope, method)

	httpFile, err := httpfile.Open(u.String())
	if err != nil {
		return nil, err
	}

	return httpFile, nil
}

func (client *HTTPClient) Call(scope string, method string, req interface{}, res interface{}) error {
	// Prepare full URL
	u, _ := url.Parse(client.url)
	u.Path = filepath.Join(u.Path, scope, method)

	// Marshal request data to JSON
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	reqReader := bytes.NewReader(reqBytes)
	httpRes, err := client.http.Post(u.String(), "application/json", reqReader)
	if err != nil {
		return err
	}

	resBytes, _ := ioutil.ReadAll(httpRes.Body)
	err = json.Unmarshal(resBytes, res)
	if err != nil {
		return err
	}

	return nil
}

func (client *HTTPClient) Open(id string) (io.ReadCloser, error) {
	resp, err := client.http.Get(client.url + id)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return resp.Body, nil
}
