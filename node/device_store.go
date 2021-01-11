package node

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/store/file"
	"github.com/cryptopunkscc/lore/util"
	"gorm.io/gorm"
	"log"
)

var _ store.Store = &DeviceStore{}

type DeviceStore struct {
	config   *Config
	db       *gorm.DB
	primary  store.Store
	observer store.Observer
	ids      id.Counter
}

func NewDeviceStore(config *Config, db *gorm.DB, storeObserver store.Observer) (*DeviceStore, error) {
	var err error

	dev := &DeviceStore{
		config:   config,
		db:       db,
		observer: storeObserver,
		ids:      id.NewCounter(),
	}

	dev.primary, err = file.NewFileStore(dev.config.GetDataDir(), dev.Added, dev.Removed)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (store *DeviceStore) Refresh() {
	primarySet, _ := store.primary.List()
	primarySet.Each(func(id id.ID) {
		store.Added(id)
	})
}

func (store *DeviceStore) Added(id id.ID) {
	if store.ids.Increment(id) == 1 {
		log.Println("<store> added", id)
		if store.observer != nil {
			store.observer.Added(id)
		}
	}
}

func (store *DeviceStore) Removed(id id.ID) {
	if store.ids.Decrement(id) == 0 {
		log.Println("<store> removed", id)
		if store.observer != nil {
			store.observer.Removed(id)
		}
	}
}

func (store *DeviceStore) Read(id id.ID) (util.ReadSeekCloser, error) {
	r, err := store.primary.Read(id)
	return r, err
}

func (store *DeviceStore) List() (id.Set, error) {
	return store.ids.Set(), nil
}

func (store *DeviceStore) Free() (uint64, error) {
	return store.primary.Free()
}

func (store *DeviceStore) Create() (store.Writer, error) {
	return store.primary.Create()
}

func (store *DeviceStore) Delete(id id.ID) error {
	return store.primary.Delete(id)
}
