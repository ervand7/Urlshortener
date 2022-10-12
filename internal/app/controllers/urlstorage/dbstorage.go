package urlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/ervand7/urlshortener/internal/app/database"
	q "github.com/ervand7/urlshortener/internal/app/database/rawqueries"
	_errors "github.com/ervand7/urlshortener/internal/app/errors"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"github.com/jackc/pgtype"
)

const (
	BufLen  int           = 3
	Timeout time.Duration = 5
)

type DBStorage struct {
	DB         database.Database
	delChan    chan string
	buf        []string
	resetTimer bool
	timer      *time.Timer
}

func NewDBStorage(db database.Database) *DBStorage {
	storage := &DBStorage{
		DB:         db,
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
	rows, err := d.DB.Conn.QueryContext(ctx, q.Set, userID, short, origin)
	if err != nil {
		return err
	}
	defer d.DB.CloseRows(rows)

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

func (d *DBStorage) SetMany(ctx context.Context, dbEntries []utils.DBEntry) error {
	transaction, err := d.DB.Conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		err = transaction.Rollback()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
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
	rows := d.DB.Conn.QueryRowContext(ctx, q.Get, short)
	var active bool

	err = rows.Scan(&origin, &active)
	if err != nil {
		return "", err
	}
	if !active {
		return "", _errors.NewURLNotActiveError(short)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	return origin, nil
}

func (d *DBStorage) GetUserURLs(
	ctx context.Context, userID string,
) (result []map[string]string, err error) {
	rows, err := d.DB.Conn.QueryContext(ctx, q.GetUserURLs, userID)
	if err != nil {
		return nil, err
	}
	defer d.DB.CloseRows(rows)

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

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *DBStorage) DeleteURLs(shortUrls []string) {
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
			utils.Logger.Error(fmt.Sprintf("delete batch: %s", err.Error()))
		}
	}()
}

func (d *DBStorage) deleteBatch(ctx context.Context, shortUrls []string) error {
	arr := &pgtype.VarcharArray{}
	if err := arr.Set(shortUrls); err != nil {
		return fmt.Errorf("prepare argument: %w", err)
	}
	res, err := d.DB.Conn.ExecContext(ctx, "UPDATE url SET active = false WHERE shorts ANY ($1)", arr)
	if err != nil {
		return fmt.Errorf("mass update: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking count: %w", err)
	}
	if count != int64(len(shortUrls)) {
		return fmt.Errorf("invalid count, expected %d, got %d", len(shortUrls), count)
	}
	return nil
}
