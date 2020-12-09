package client

import (
	"errors"
	"github.com/cryptopunkscc/lore/comm/proto"
)

const scopeAdmin = "admin"

type AdminClient struct {
	api APIClient
}

func (admin *AdminClient) Add(path string) (string, error) {
	var err error
	var response proto.AdminAddResponse
	var request = proto.AdminAddRequest{
		Path: path,
	}

	err = admin.api.Call("admin", "add", &request, &response)
	if err != nil {
		return "", err
	}

	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.ID, nil
}

func (admin *AdminClient) List() ([]string, error) {
	var err error
	var response proto.AdminListResponse

	err = admin.api.Call(scopeAdmin, "list", &proto.AdminListRequest{}, &response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Items, nil
}

func (admin *AdminClient) ListSources() ([]string, error) {
	var err error
	var response proto.AdminListSourcesResponse

	err = admin.api.Call("admin", "listsources", &proto.AdminListSourcesRequest{}, &response)
	if err != nil {
		return nil, err
	}

	return response.Sources, nil
}

func (admin *AdminClient) AddSource(address string) error {
	var err error
	var response proto.AdminAddSourceResponse

	err = admin.api.Call(scopeAdmin, "addsource", &proto.AdminAddSourceRequest{
		Address: address,
	}, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}

func (admin *AdminClient) RemoveSource(address string) error {
	var err error
	var response proto.AdminRemoveSourceResponse

	err = admin.api.Call(scopeAdmin, "removesource", &proto.AdminRemoveSourceRequest{
		Address: address,
	}, &response)
	if err != nil {
		return err
	}
	if response.Error != "" {
		return errors.New(response.Error)
	}
	return nil
}
