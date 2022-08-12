package models

import "sync"

type URLRepository interface {
	Get(short string) (origin string)
	Set(short, origin string)
}

type URLStorage struct {
	HashTable map[string]string
	Mutex     sync.Mutex
}

func (s *URLStorage) Get(short string) (origin string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	origin, exist := s.HashTable[short]
	if !exist {
		return ""
	}
	return origin
}

func (s *URLStorage) Set(short, origin string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.HashTable[short] = origin
}
