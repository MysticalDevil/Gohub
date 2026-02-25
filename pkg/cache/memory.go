package cache

import (
	"strconv"
	"sync"
	"time"
)

type memoryItem struct {
	value     string
	expiresAt time.Time
	hasExpire bool
}

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]memoryItem),
	}
}

func (store *MemoryStore) Set(key, value string, expireTime time.Duration) {
	store.mu.Lock()
	defer store.mu.Unlock()

	item := memoryItem{value: value}
	if expireTime > 0 {
		item.hasExpire = true
		item.expiresAt = time.Now().Add(expireTime)
	}
	store.items[key] = item
}

func (store *MemoryStore) Get(key string) string {
	store.mu.RLock()
	item, ok := store.items[key]
	store.mu.RUnlock()

	if !ok {
		return ""
	}

	if item.hasExpire && time.Now().After(item.expiresAt) {
		store.mu.Lock()
		delete(store.items, key)
		store.mu.Unlock()
		return ""
	}

	return item.value
}

func (store *MemoryStore) Has(key string) bool {
	return store.Get(key) != ""
}

func (store *MemoryStore) Forget(key string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.items, key)
}

func (store *MemoryStore) Forever(key, value string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.items[key] = memoryItem{value: value}
}

func (store *MemoryStore) Flush() {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.items = make(map[string]memoryItem)
}

func (store *MemoryStore) IsAlive() error {
	return nil
}

func (store *MemoryStore) Increment(parameters ...any) {
	store.adjustInt(true, parameters...)
}

func (store *MemoryStore) Decrement(parameters ...any) {
	store.adjustInt(false, parameters...)
}

func (store *MemoryStore) adjustInt(increment bool, parameters ...any) {
	if len(parameters) == 0 {
		return
	}

	key, _ := parameters[0].(string)
	if key == "" {
		return
	}

	var delta int64 = 1
	if len(parameters) > 1 {
		if v, ok := parameters[1].(int64); ok {
			delta = v
		}
	}

	value, _ := strconv.ParseInt(store.Get(key), 10, 64)
	if increment {
		value += delta
	} else {
		value -= delta
	}
	store.Set(key, strconv.FormatInt(value, 10), 0)
}
