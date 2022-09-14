package database

import (
	"database/sql"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/utils"
	q "github.com/ervand7/urlshortener/internal/app/utils/rawqueries"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

var DB Database

type Database struct {
	Conn *sql.DB
}

func (d *Database) ConnStart() (err error) {
	db, err := sql.Open("pgx",
		config.GetConfig().DatabaseDSN)
	if err != nil {
		return err
	}
	d.Conn = db
	return nil
}

func (d *Database) ConnClose() (err error) {
	err = d.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) SetConnPool() {
	d.Conn.SetMaxOpenConns(20)
	d.Conn.SetMaxIdleConns(20)
	d.Conn.SetConnMaxIdleTime(time.Second * 30)
	d.Conn.SetConnMaxLifetime(time.Minute * 2)
}

func (d *Database) Ping() error {
	err := d.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CreateAll() error {
	_, err := d.Conn.Exec(q.CreateAll)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		utils.Logger.Error(err.Error())
	}
}
