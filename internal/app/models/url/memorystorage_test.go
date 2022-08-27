package url

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	memoryStorage := MemoryStorage{
		HashTable: make(map[string]string, 0),
	}

	key := "Hello"
	value := "World"
	memoryStorage.Set(key, value)
	assert.Equal(t, memoryStorage.HashTable[key], value)

	result := memoryStorage.Get(key)
	assert.Equal(t, result, value)
}
