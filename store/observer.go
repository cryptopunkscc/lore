package store

type EventFunc func(id string)

type Observer interface {
	Added(id string)
	Removed(id string)
}
