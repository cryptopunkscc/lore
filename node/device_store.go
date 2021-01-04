package node

import (
	"github.com/cryptopunkscc/lore/store"
	"gorm.io/gorm"
	"log"
)

var _ store.Store = &DeviceStore{}

type DeviceStore struct {
	config         *Config
	fileMapRepo    store.FileMapRepo
	primary        store.Store
	readableStores store.Group
	networkStores  store.Group
	db             *gorm.DB

	storeObserver store.Observer
}

func (dev *DeviceStore) Added(id string) {
	log.Println(id, "added")
	if dev.storeObserver != nil {
		dev.storeObserver.Added(id)
	}
}

func (dev *DeviceStore) Removed(id string) {
	log.Println(id, "removed")
	if dev.storeObserver != nil {
		dev.storeObserver.Removed(id)
	}
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

	dev.primary, err = store.NewFileStore(dev.config.GetDataDir())
	if err != nil {
		return nil, err
	}

	_ = dev.readableStores.Add(dev.primary)

	dev.fileMapRepo, err = store.NewFileMapRepoGorm(dev.db)
	if err != nil {
		return nil, err
	}

	return dev, nil
}

func (dev *DeviceStore) AddLocalDir(path string) error {
	m, _ := store.NewFileMapStore(dev.fileMapRepo)

	err := dev.readableStores.Add(m)
	if err != nil {
		return err
	}

	err = store.NewDirSync(path, m, dev).Update()
	if err != nil {
		return err
	}

	return nil
}

func (dev *DeviceStore) AddNetworkStore(url string) error {
	s := store.NewHTTPStore(url)
	return dev.networkStores.Add(s)
}

func (dev *DeviceStore) Read(id string) (store.ReadSeekCloser, error) {
	r, err := dev.readableStores.Read(id)
	if err != nil {
		return dev.networkStores.Read(id)
	}
	return r, err
}

func (dev *DeviceStore) List() ([]string, error) {
	return dev.readableStores.List()
}

func (dev *DeviceStore) Free() (int64, error) {
	return dev.primary.Free()
}

func (dev *DeviceStore) Create() (store.Writer, error) {
	w, err := dev.primary.Create()
	if err != nil {
		return nil, err
	}

	ww := store.NewWrappedWriter(w, func(id string, err error) error {
		if err == nil {
			dev.Added(id)
		}
		return err
	})
	return ww, nil
}

func (dev *DeviceStore) Delete(id string) error {
	err := dev.primary.Delete(id)
	if err == nil {
		dev.Removed(id)
	}
	return err
}
