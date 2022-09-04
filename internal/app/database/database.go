package database

import (
	"database/sql"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/utils"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Database struct {
	Conn *sql.DB
}

func (d *Database) StartConn() (err error) {
	db, err := sql.Open("pgx",
		config.GetConfig().DatabaseDSN)
	if err != nil {
		return err
	}

	d.Conn = db
	return nil
}

func (d *Database) Close() {
	err := d.Conn.Close()
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}

func (d *Database) Ping() error {
	err := d.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}
