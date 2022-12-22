package storage

import (
	"context"
	"github.com/ervand7/urlshortener/internal/database"
	_errors "github.com/ervand7/urlshortener/internal/errors"
	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
	"github.com/lib/pq"
	"time"
)

const (
	BufLen   int           = 3
	Timeout  time.Duration = 2
	setQuery string        = `
		with cte as (
			insert into url ("user_id", "short", "origin")
				values ($1, $2, $3)
				on conflict ("origin") do nothing
				returning "short")
		select 'null'
		where exists(select 1 from cte)
		union all
		select "short"
		from url
		where "origin" = $3
		  and not exists(select 1 from cte);
`
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
	rows, err := d.db.Conn.QueryContext(ctx, setQuery, userID, short, origin)
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

func (d *DBStorage) SetMany(ctx context.Context, dbEntries []models.Entry) error {
	transaction, err := d.db.Conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		transaction.Rollback()
	}()

	stmt, err := transaction.PrepareContext(ctx, setQuery)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			logger.Logger.Error(err.Error())
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
	row := d.db.Conn.QueryRowContext(
		ctx, `select "origin", "active" from url where "short" = $1;`, short,
	)
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
	rows, err := d.db.Conn.QueryContext(
		ctx, `select "short", "origin" from url where "user_id" = $1;`, userID,
	)
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
			logger.Logger.Error(err.Error())
		}
	}()
}

func (d *DBStorage) deleteBatch(ctx context.Context, shortUrls []string) error {
	_, err := d.db.Conn.ExecContext(
		ctx,
		`update url set "active" = false  where "short" = ANY($1)`,
		pq.Array(shortUrls),
	)
	if err != nil {
		return err
	}

	return nil
}
