package client

type Client struct {
	api APIClient
}

func NewClient(address string) *Client {
	c := &Client{}
	c.api = NewHTTPClient(address)
	return c
}

func (client *Client) Admin() *AdminClient {
	return &AdminClient{api: client.api}
}

func (client *Client) Item() *ItemClient {
	return &ItemClient{api: client.api}
}

func (client *Client) Download() *DownloadClient {
	return &DownloadClient{api: client.api}
}

func (client *Client) Local() *LocalClient {
	return &LocalClient{api: client.api}
}
