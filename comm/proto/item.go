package proto

type ItemInfoRequest struct {
	ID string
}

type ItemInfoResponse struct {
	ID    string
	Error string
}
