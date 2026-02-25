package captcha

import "sync"

// MemoryStore implements base64Captcha.Store for tests.
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]string),
	}
}

func (s *MemoryStore) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
	return nil
}

func (s *MemoryStore) Get(key string, clear bool) string {
	s.mu.RLock()
	val, ok := s.items[key]
	s.mu.RUnlock()
	if !ok {
		return ""
	}
	if clear {
		s.mu.Lock()
		delete(s.items, key)
		s.mu.Unlock()
	}
	return val
}

func (s *MemoryStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
