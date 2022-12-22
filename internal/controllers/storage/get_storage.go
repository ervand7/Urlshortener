package storage

import (
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/storage/db_storage"
	f "github.com/ervand7/urlshortener/internal/controllers/storage/file_storage"
	m "github.com/ervand7/urlshortener/internal/controllers/storage/memory_storage"
	"github.com/ervand7/urlshortener/internal/models"
)

func GetStorage() (storage models.Storage) {
	if config.GetConfig().DatabaseDSN != "" {
		db := db_storage.Database{}
		db.Run()
		storage = db_storage.NewDBStorage(db)
		return storage
	}
	if config.GetConfig().FileStoragePath != "" {
		storage = &f.FileStorage{}
		return storage
	}
	storage = &m.MemoryStorage{
		HashTable: make(map[string]m.ShortenURLStruct, 0),
	}
	return storage
}
