package storage

type Storage interface {
	Create() (Writer, error)
	Delete(id string) error
	Open(id string) (ReadSeekCloser, error)
	List() ([]string, error)
	Contains(id string) (bool, error)
}
