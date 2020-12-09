package client

import "io"

type APIClient interface {
	Call(scope string, method string, req interface{}, res interface{}) error
	GetStream(scope string, method string) (io.Reader, error)
}
