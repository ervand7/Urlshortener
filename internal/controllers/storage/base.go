package storage

import (
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/storage/filestorage"
	"github.com/ervand7/urlshortener/internal/database"
	"github.com/ervand7/urlshortener/internal/models"
)

func GetStorage() (storage models.Storage) {
	if config.GetConfig().DatabaseDSN != "" {
		db := database.Database{}
		db.Run()
		storage = NewDBStorage(db)
		return storage
	}
	if config.GetConfig().FileStoragePath != "" {
		storage = &filestorage.FileStorage{}
		return storage
	}
	storage = &MemoryStorage{
		HashTable: make(map[string]ShortenURLStruct, 0),
	}
	return storage
}
