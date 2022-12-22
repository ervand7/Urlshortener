package dbstorage

import (
	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/pressly/goose/v3"
)

// Downgrade for test DB
func Downgrade() {
	db := Database{}
	db.Run()
	defer func() {
		if err := goose.Run("down", db.Conn, getMigrationsDir()); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()
}
