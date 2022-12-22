package memorystorage

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStorage(t *testing.T) {
	memoryStorage := MemoryStorage{
		HashTable: make(map[string]ShortenURLStruct, 0),
	}

	short := "Hello"
	origin := "World"
	userID := uuid.New().String()
	err := memoryStorage.Set(context.TODO(), userID, short, origin)
	assert.NoError(t, err)
	assert.Equal(t, memoryStorage.HashTable[short].Origin, origin)

	result, err := memoryStorage.Get(context.TODO(), short)
	assert.NoError(t, err)
	assert.Equal(t, result, origin)
}
