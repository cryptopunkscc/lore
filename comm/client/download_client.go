package client

import "io"

type DownloadClient struct {
	api APIClient
}

func (c *DownloadClient) Stream(id string) (io.Reader, error) {
	return c.api.GetStream("download", id)
}
