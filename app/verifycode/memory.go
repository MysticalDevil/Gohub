package verifycode

import "sync"

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]string),
	}
}

func (store *MemoryStore) Set(id string, value string) bool {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.items[id] = value
	return true
}

func (store *MemoryStore) Get(id string, clear bool) string {
	store.mu.RLock()
	value, ok := store.items[id]
	store.mu.RUnlock()

	if !ok {
		return ""
	}

	if clear {
		store.mu.Lock()
		delete(store.items, id)
		store.mu.Unlock()
	}

	return value
}

func (store *MemoryStore) Verify(id, answer string, clear bool) bool {
	value := store.Get(id, clear)
	return value == answer
}
