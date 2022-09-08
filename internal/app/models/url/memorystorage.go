package url

import "sync"

type MemoryStorage struct {
	HashTable map[string]ShortenURLStruct
	Mutex     sync.Mutex
}

type ShortenURLStruct struct {
	Short  string
	Origin string
	UserID string
}

func (m *MemoryStorage) Set(userID, short, origin string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HashTable[short] = ShortenURLStruct{
		Short:  short,
		Origin: origin,
		UserID: userID,
	}
}

func (m *MemoryStorage) Get(short string) (origin string) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	origin = m.HashTable[short].Origin
	return origin
}

func (m *MemoryStorage) GetUserURLs(userID string) (userURLs []map[string]string, err error) {
	userURLs = make([]map[string]string, 0)
	for _, data := range m.HashTable {
		if data.UserID == userID {
			pair := map[string]string{
				"short_url":    data.Short,
				"original_url": data.Origin,
			}
			userURLs = append(userURLs, pair)
		}
	}

	return userURLs, nil
}
