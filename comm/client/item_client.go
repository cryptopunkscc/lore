package client

import (
	"errors"
	"github.com/cryptopunkscc/lore/comm/proto"
)

type ItemClient struct {
	api APIClient
}

func (item *ItemClient) Info(id string) (bool, error) {
	var err error
	var response proto.ItemInfoResponse

	err = item.api.Call("item", "info", &proto.ItemInfoRequest{ID: id}, &response)
	if err != nil {
		return false, err
	}

	if response.Error != "" {
		return false, errors.New(response.Error)
	}

	return true, nil
}
