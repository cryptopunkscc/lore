package storage

type Writer interface {
	Write(data []byte) (int, error)
	Finalize() (string, error)
	Discard() error
}
