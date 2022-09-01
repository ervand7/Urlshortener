package url

import "sync"

type ShortenURLStruct struct {
	Shorten string
	Origin  string
	UserID  string
}

type MemoryStorage struct {
	HashTable map[string]ShortenURLStruct
	Mutex     sync.Mutex
}

func (m *MemoryStorage) Get(short string) (origin string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	origin = m.HashTable[short].Origin
	return origin
}

func (m *MemoryStorage) Set(userID, short, origin string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HashTable[short] = ShortenURLStruct{
		Shorten: short,
		Origin:  origin,
		UserID:  userID,
	}
}
