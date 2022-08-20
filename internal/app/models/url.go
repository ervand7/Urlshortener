package models

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"sync"
)

type URLRepository interface {
	Get(short string) (origin string)
	Set(short, origin string)
}

type URLStorage struct {
	HashTable map[string]string
	FileTable URLRepository
	Mutex     sync.Mutex
}

func (s *URLStorage) Get(short string) (origin string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	switch config.GetConfig().FileStoragePath {
	case "":
		origin, exist := s.HashTable[short]
		if !exist {
			return ""
		}
		return origin

	default:
		origin := s.FileTable.Get(short)
		return origin

	}
}

func (s *URLStorage) Set(short, origin string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	switch config.GetConfig().FileStoragePath {
	case "":
		s.HashTable[short] = origin
	default:
		s.FileTable.Set(short, origin)
	}
}
