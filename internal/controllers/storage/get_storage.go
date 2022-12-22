package storage

import (
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
	f "github.com/ervand7/urlshortener/internal/controllers/storage/filestorage"
	m "github.com/ervand7/urlshortener/internal/controllers/storage/memorystorage"
	"github.com/ervand7/urlshortener/internal/models"
)

func GetStorage() (storage models.Storage) {
	if config.GetConfig().DatabaseDSN != "" {
		db := dbstorage.Database{}
		db.Run()
		storage = dbstorage.NewDBStorage(db)
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