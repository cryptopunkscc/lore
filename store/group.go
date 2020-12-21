package store

type Group interface {
	Store
	Add(store Store) error
	Remove(store Store) error
}
