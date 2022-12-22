package dbstorage

import (
	"github.com/pressly/goose/v3"

	"github.com/ervand7/urlshortener/internal/logger"
)

// Downgrade for test DB.
func Downgrade() {
	db := Database{}
	db.Run()
	defer func() {
		if err := goose.Run("down", db.Conn, getMigrationsDir()); err != nil {
			logger.Logger.Error(err.Error())
		}
	}()
}
