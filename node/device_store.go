package node

import (
	"github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store"
	"github.com/cryptopunkscc/lore/store/pathcache"
	"gorm.io/gorm"
	"log"
)

var _ store.Store = &DeviceStore{}

type DeviceStore struct {
	config         *Config
	fileMapRepo    pathcache.Repo
	primary        store.Store
	readableStores store.Group
	networkStores  store.Group
	db             *gorm.DB

	storeObserver store.Observer
}

func NewDeviceStore(config *Config, db *gorm.DB, storeObserver store.Observer) (*DeviceStore, error) {
	var err error

	dev := &DeviceStore{
		config:         config,
		readableStores: store.NewSeqGroup(),
		networkStores:  store.NewSeqGroup(),
		db:             db,
		storeObserver:  storeObserver,
	}

	dev.primary, err = store.NewFileStore(dev.config.GetDataDir(), dev.Added, dev.Removed)
	if err != nil {
		return nil, err
	}

	_ = dev.readableStores.Add(dev.primary)

	return dev, nil
}

func (dev *DeviceStore) Added(id id.ID) {
	log.Println("<store> added", id)
	if dev.storeObserver != nil {
		dev.storeObserver.Added(id)
	}
}

func (dev *DeviceStore) Removed(id id.ID) {
	log.Println("<store> removed", id)
	if dev.storeObserver != nil {
		dev.storeObserver.Removed(id)
	}
}

func (dev *DeviceStore) AddNetworkStore(url string) error {
	s := store.NewHTTPStore(url)
	return dev.networkStores.Add(s)
}

func (dev *DeviceStore) Read(id id.ID) (store.ReadSeekCloser, error) {
	r, err := dev.readableStores.Read(id)
	if err != nil {
		return dev.networkStores.Read(id)
	}
	return r, err
}

func (dev *DeviceStore) List() (id.Set, error) {
	return dev.readableStores.List()
}

func (dev *DeviceStore) Free() (uint64, error) {
	return dev.primary.Free()
}

func (dev *DeviceStore) Create() (store.Writer, error) {
	return dev.primary.Create()
}

func (dev *DeviceStore) Delete(id id.ID) error {
	return dev.primary.Delete(id)
}
