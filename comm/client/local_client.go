package client

import "io"

type LocalClient struct {
	api APIClient
}

func (c *LocalClient) Stream(id string) (io.Reader, error) {
	return c.api.GetStream("local", id)
}
