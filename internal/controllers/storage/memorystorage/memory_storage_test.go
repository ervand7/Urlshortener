package memorystorage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStorageGetSet(t *testing.T) {
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

func TestMemoryStorageGetUserURLs(t *testing.T) {
	memoryStorage := MemoryStorage{
		HashTable: make(map[string]ShortenURLStruct, 0),
	}
	userID := uuid.New().String()

	short1 := "Hello1"
	origin1 := "World1"
	err := memoryStorage.Set(context.TODO(), userID, short1, origin1)
	assert.NoError(t, err)

	short2 := "Hello2"
	origin2 := "World2"
	err = memoryStorage.Set(context.TODO(), userID, short2, origin2)
	assert.NoError(t, err)

	result, err := memoryStorage.GetUserURLs(context.TODO(), userID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	for _, val := range result {
		_, exists := val["original_url"]
		assert.True(t, exists)
		_, exists = val["short_url"]
		assert.True(t, exists)
	}
}
