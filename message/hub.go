package message

import (
	"sync"
)

type Hub struct {
	stores       map[string]*Store
	store_add    chan *Store
	store_remove chan *Store
	ch_backup    chan *ChMsg
	rwlock       sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		stores:       make(map[string]*Store),
		store_add:    make(chan *Store),
		store_remove: make(chan *Store),
		ch_backup:    make(chan *ChMsg, 20),
	}
}

func (h *Hub) AddStore(store *Store) {
	h.rwlock.Lock()
	defer h.rwlock.Unlock()
	h.stores[store.id] = store
}

func (h *Hub) GetStore(store_id string) (*Store, bool) {
	h.rwlock.RLock()
	defer h.rwlock.RUnlock()
	store, ok := h.stores[store_id]
	return store, ok
}

func (h *Hub) RemoveStore(store *Store) {
	close(store.register)
	close(store.unregister)
	close(store.transfer)
	h.rwlock.Lock()
	defer h.rwlock.Unlock()
	delete(h.stores, store.id)
}

func (h *Hub) BackupMsg(msg *ChMsg) {
	logger.Info("Backup Msg", msg.DataStr)
}

func (h *Hub) run() {
	for {
		select {
		case store := <-h.store_add:
			h.AddStore(store)
		case store := <-h.store_remove:
			h.RemoveStore(store)
		case msg := <-h.ch_backup:
			go h.BackupMsg(msg)
		}
	}
}
