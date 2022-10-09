package storage

import (
	"context"
	"github.com/ervand7/urlshortener/internal/app/database"
	q "github.com/ervand7/urlshortener/internal/app/database/rawqueries"
	_errors "github.com/ervand7/urlshortener/internal/app/errors"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/lib/pq"
	"time"
)

const (
	BufLen  int           = 3
	Timeout time.Duration = 2
)

type DBStorage struct {
	db         database.Database
	delChan    chan string
	buf        []string
	resetTimer bool
	timer      *time.Timer
}

func NewDBStorage(db database.Database) *DBStorage {
	storage := &DBStorage{
		db:         db,
		delChan:    make(chan string),
		buf:        make([]string, 0, BufLen),
		resetTimer: true,
		timer:      time.NewTimer(0),
	}
	go storage.startWorker()
	return storage
}

func (d *DBStorage) Set(ctx context.Context, userID, short, origin string) error {
	var existsShort string
	rows, err := d.db.Conn.QueryContext(ctx, q.Set, userID, short, origin)
	if err != nil {
		return err
	}
	defer d.db.CloseRows(rows)

	for rows.Next() {
		err = rows.Scan(&existsShort)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	if existsShort != "null" {
		return _errors.NewShortAlreadyExistsError(existsShort)
	}

	return nil
}

func (d *DBStorage) SetMany(ctx context.Context, dbEntries []utils.Entry) error {
	transaction, err := d.db.Conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		transaction.Rollback()
	}()

	stmt, err := transaction.PrepareContext(ctx, q.Set)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}()

	for _, e := range dbEntries {
		if _, err = stmt.ExecContext(ctx, e.UserID, e.Short, e.Origin); err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func (d *DBStorage) Get(ctx context.Context, short string) (origin string, err error) {
	row := d.db.Conn.QueryRowContext(ctx, q.Get, short)
	var active bool

	err = row.Scan(&origin, &active)
	if err != nil {
		return "", err
	}
	if !active {
		return "", _errors.NewURLNotActiveError(short)
	}

	return origin, nil
}

func (d *DBStorage) GetUserURLs(
	ctx context.Context, userID string,
) (result []map[string]string, err error) {
	rows, err := d.db.Conn.QueryContext(ctx, q.GetUserURLs, userID)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer d.db.CloseRows(rows)

	var (
		short  string
		origin string
	)
	result = make([]map[string]string, 0)
	for rows.Next() {
		err = rows.Scan(&short, &origin)
		if err != nil {
			return nil, err
		}
		pair := map[string]string{
			"short_url":    short,
			"original_url": origin,
		}
		result = append(result, pair)
	}

	return result, nil
}

func (d *DBStorage) DeleteUserURLs(shortUrls []string) {
	for _, val := range shortUrls {
		go func(short string) {
			d.delChan <- short
		}(val)
	}
}

func (d *DBStorage) startWorker() {
	ctx := context.Background()
	for {
		select {
		case shortURL := <-d.delChan:
			if d.resetTimer {
				d.timer.Reset(time.Second * Timeout)
				d.resetTimer = false
			}
			d.buf = append(d.buf, shortURL)
			if len(d.buf) >= BufLen {
				d.flush(ctx)
				d.resetTimer = true
			}
		case <-d.timer.C:
			if len(d.buf) > 0 {
				d.flush(ctx)
			}
		}
	}
}

func (d *DBStorage) flush(ctx context.Context) {
	toDelete := make([]string, len(d.buf))
	copy(toDelete, d.buf)
	d.buf = make([]string, 0)
	go func() {
		err := d.deleteBatch(ctx, toDelete)
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}()
}

func (d *DBStorage) deleteBatch(ctx context.Context, shortUrls []string) error {
	_, err := d.db.Conn.ExecContext(ctx, q.DeleteURL, pq.Array(shortUrls))
	if err != nil {
		return err
	}

	return nil
}
