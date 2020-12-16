package store

type Group interface {
	ReadEditor
	Add(store ReadEditor) error
	Remove(store ReadEditor) error
}
