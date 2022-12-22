package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
)

// MemoryStorage implementation of Storage interface for working with memory
type MemoryStorage struct {
	HashTable map[string]ShortenURLStruct
	Mutex     sync.Mutex
}

// ShortenURLStruct struct for storing data in memory
type ShortenURLStruct struct {
	Short  string
	Origin string
	UserID string
}

// Set write one entity
func (m *MemoryStorage) Set(
	_ context.Context, userID, short, origin string,
) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	m.HashTable[short] = ShortenURLStruct{
		Short:  short,
		Origin: origin,
		UserID: userID,
	}
	return nil
}

// Get one entity
func (m *MemoryStorage) Get(
	_ context.Context, short string,
) (origin string, err error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
	origin = m.HashTable[short].Origin
	return origin, nil
}

// GetUserURLs gets all user urls
func (m *MemoryStorage) GetUserURLs(
	_ context.Context, userID string,
) (userURLs []map[string]string, err error) {
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

// SetMany batch write
func (m *MemoryStorage) SetMany(_ context.Context, _ []models.Entry) error {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return err
}

// DeleteUserURLs batch delete
func (m *MemoryStorage) DeleteUserURLs(_ []string) {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
}
