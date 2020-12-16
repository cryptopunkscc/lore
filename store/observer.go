package store

type Observer interface {
	Added(id string)
	Removed(id string)
}
