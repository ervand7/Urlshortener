package database

import (
	"context"
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/jackc/pgx/v4"
)

type Database struct {
	Conn *pgx.Conn
}

func (d *Database) StartConn() (err error) {
	d.Conn, err = pgx.Connect(context.Background(), config.GetConfig().DatabaseDSN)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Close() {
	err := d.Conn.Close(context.Background())
	if err != nil {
		utils.Logger.Error(err.Error())
	}
}

func (d *Database) Ping() error {
	err := d.Conn.Ping(context.Background())
	if err != nil {
		return err
	}
	return nil
}
