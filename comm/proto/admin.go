package proto

// admin/add - add a local file to storage
type AdminAddRequest struct {
	Path string `json:"path"`
}

type AdminAddResponse struct {
	ID      string `json:"id"`
	Size    int    `json:"size"`
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Error   string
}

// admin/list - list all files in the storage
type AdminListRequest struct {
}

type AdminListResponse struct {
	Items []string
	Error string
}

// admin/addsource - add a URL to the sources list
type AdminAddSourceRequest struct {
	Address string
}

type AdminAddSourceResponse struct {
	Error string
}

// admin/removesource
type AdminRemoveSourceRequest struct {
	Address string
}

type AdminRemoveSourceResponse struct {
	Error string
}

// admin/liststources - list all URLs on the sources list
type AdminListSourcesRequest struct{}

type AdminListSourcesResponse struct {
	Sources []string
}
