package dbstorage

import (
	"database/sql"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	maxOpenConnections    = 20
	maxIdleConnections    = 20
	connMaxIdleTimeSecond = 30
	connMaxLifetimeSecond = 2
)

// Database for manage DB settings
type Database struct {
	Conn *sql.DB
}

// Run launches DB
func (db *Database) Run() {
	err := db.ConnStart()
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
	db.setConnPool()
	err = db.migrate()
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}
}

// ConnStart opens DB connection
func (db *Database) ConnStart() (err error) {
	conn, err := goose.OpenDBWithDriver("pgx", config.GetDatabaseDSN())
	if err != nil {
		return err
	}
	db.Conn = conn
	return nil
}

// ConnClose closes DB connection
func (db *Database) ConnClose() (err error) {
	err = db.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Ping is used for check if DB connection is alive
func (db *Database) Ping() (err error) {
	err = db.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

// CloseRows closes db.Conn.Query rows
func (db *Database) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		logger.Logger.Error(err.Error())
	}
}

func (db *Database) migrate() (err error) {
	if err = goose.Run("up", db.Conn, getMigrationsDir()); err != nil {
		return err
	}
	return nil
}

func (db *Database) setConnPool() {
	db.Conn.SetMaxOpenConns(maxOpenConnections)
	db.Conn.SetMaxIdleConns(maxIdleConnections)
	db.Conn.SetConnMaxIdleTime(time.Second * connMaxIdleTimeSecond)
	db.Conn.SetConnMaxLifetime(time.Minute * connMaxLifetimeSecond)
}

func getMigrationsDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	migrationsDir := filepath.Join(currentDir, "/../../../../migrations")
	return migrationsDir
}
