package url

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	memoryStorage := MemoryStorage{
		HashTable: make(map[string]ShortenURLStruct, 0),
	}

	key := "Hello"
	value := "World"
	userID := uuid.New().String()
	memoryStorage.Set(userID, key, value)
	assert.Equal(t, memoryStorage.HashTable[key].Origin, value)

	result := memoryStorage.Get(key)
	assert.Equal(t, result, value)
}
