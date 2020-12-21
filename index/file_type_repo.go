package index

type FileTypeRepo interface {
	Set(id string, typ string) error
	Get(id string) string
	Clear(id string) error
}
