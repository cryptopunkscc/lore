package store

import (
	"errors"
	_id "github.com/cryptopunkscc/lore/id"
)

var _ Group = &SeqGroup{}

// Group provides a Store interface over a collection of stores.
type SeqGroup struct {
	stores []Store
}

// NewSeqGroup makes a new instance of a SeqGroup
func NewSeqGroup() *SeqGroup {
	m := &SeqGroup{}
	m.stores = make([]Store, 0)
	return m
}

// Add adds a store to the collection
func (group *SeqGroup) Add(store Store) error {
	group.stores = append(group.stores, store)
	return nil
}

// Remove removes a store from the collection
func (group *SeqGroup) Remove(store Store) error {
	for i, s := range group.stores {
		if s == store {
			group.stores = append(group.stores[:i], group.stores[i+1:]...)
		}
	}
	return nil
}

// Read will call Read on every store in the collection and return the result of the first successful call.
func (group *SeqGroup) Read(id _id.ID) (ReadSeekCloser, error) {
	for _, s := range group.stores {
		r, err := s.Read(id)
		if err == nil {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

// List returns a merged list of files from all stores
func (group *SeqGroup) List() (_id.Set, error) {
	set := _id.NewSet()

	for _, store := range group.stores {
		subset, err := store.List()
		if err != nil {
			continue
		}

		subset.Each(func(id _id.ID) {
			set.Add(id)
		})
	}

	return set, nil
}

func (group *SeqGroup) Free() (uint64, error) {
	return 0, ErrUnsupported
}

// Create will call Create on every store in the collection and return the result of the first successful call.
func (group *SeqGroup) Create() (Writer, error) {
	for _, s := range group.stores {
		w, err := s.Create()
		if err == nil {
			return w, nil
		}
	}
	return nil, errors.New("create failed in every store")
}

// TODO: Delete is not yet supported.
func (group *SeqGroup) Delete(_id.ID) error {
	return ErrUnsupported
}
