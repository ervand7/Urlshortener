package storage

import (
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage"
	f "github.com/ervand7/urlshortener/internal/controllers/storage/filestorage"
	m "github.com/ervand7/urlshortener/internal/controllers/storage/memorystorage"
	"github.com/ervand7/urlshortener/internal/models"
)

// GetStorage returns Storage depending on env settings
func GetStorage() (storage models.Storage) {
	if config.GetDatabaseDSN() != "" {
		db := dbstorage.Database{}
		db.Run()
		storage = dbstorage.NewDBStorage(db)
		return storage
	}
	if config.GetFileStoragePath() != "" {
		storage = &f.FileStorage{}
		return storage
	}
	storage = &m.MemoryStorage{
		HashTable: make(map[string]m.ShortenURLStruct, 0),
	}
	return storage
}
