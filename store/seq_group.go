package store

import "errors"

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
func (group *SeqGroup) Read(id string) (ReadSeekCloser, error) {
	for _, s := range group.stores {
		r, err := s.Read(id)
		if err == nil {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

// List returns a merged list of files from all stores
func (group *SeqGroup) List() ([]string, error) {
	res := make([]string, 0)
	ids := make(map[string]bool)

	for _, store := range group.stores {
		list, err := store.List()
		if err == nil {
			for _, i := range list {
				if _, ok := ids[i]; !ok {
					res = append(res, i)
					ids[i] = true
				}
			}
		}
	}

	return res, nil
}

func (group *SeqGroup) Free() (int64, error) {
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
func (group *SeqGroup) Delete(id string) error {
	return ErrUnsupported
}
