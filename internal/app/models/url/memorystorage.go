package url

import "sync"

type MemoryStorage struct {
	HashTable map[string]string
	Mutex     sync.Mutex
}

func (m *MemoryStorage) Get(short string) (origin string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	origin, exist := m.HashTable[short]
	if !exist {
		return ""
	}
	return origin
}

func (m *MemoryStorage) Set(short, origin string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HashTable[short] = origin
}
