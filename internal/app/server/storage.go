package server

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	u "github.com/ervand7/urlshortener/internal/app/controllers/urlstorage"
	"github.com/ervand7/urlshortener/internal/app/controllers/urlstorage/filestorage"
	"github.com/ervand7/urlshortener/internal/app/database"
	"github.com/ervand7/urlshortener/internal/app/models"
)

func GetServerStorage() (storage models.Storage) {
	if config.GetConfig().DatabaseDSN != "" {
		db := database.Database{}
		db.Run()
		storage = u.NewDBStorage(db)
		return storage
	}
	if config.GetConfig().FileStoragePath != "" {
		storage = &filestorage.FileStorage{}
		return storage
	}
	storage = &u.MemoryStorage{
		HashTable: make(map[string]u.ShortenURLStruct, 0),
	}
	return storage
}
