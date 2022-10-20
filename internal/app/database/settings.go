package database

import (
	"database/sql"
	"github.com/ervand7/urlshortener/internal/app/config"
	q "github.com/ervand7/urlshortener/internal/app/database/rawqueries"
	"github.com/ervand7/urlshortener/internal/app/utils"
	_ "github.com/jackc/pgx/v4/stdlib"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Database struct {
	Conn *sql.DB
}

func (db *Database) Run() {
	err := db.ConnStart()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}
	db.SetConnPool()
	err = db.CreateAll()
	if err != nil {
		utils.Logger.Fatal(err.Error())
	}

	ch := make(chan os.Signal, 3)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		signal.Stop(ch)
		err = db.ConnClose()
		if err != nil {
			utils.Logger.Error(err.Error())
		} else {
			utils.Logger.Info("Connection to DB was closed")
		}
		os.Exit(0)
	}()
}

func (db *Database) ConnStart() (err error) {
	conn, err := sql.Open("pgx",
		config.GetConfig().DatabaseDSN)
	if err != nil {
		return err
	}
	db.Conn = conn
	return nil
}

func (db *Database) ConnClose() (err error) {
	err = db.Conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) SetConnPool() {
	db.Conn.SetMaxOpenConns(20)
	db.Conn.SetMaxIdleConns(20)
	db.Conn.SetConnMaxIdleTime(time.Second * 30)
	db.Conn.SetConnMaxLifetime(time.Minute * 2)
}

func (db *Database) Ping() error {
	err := db.Conn.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) CreateAll() error {
	_, err := db.Conn.Exec(q.CreateAll)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		utils.Logger.Error(err.Error())
	}
}
