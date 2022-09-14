package url

import (
	"context"
	"github.com/ervand7/urlshortener/internal/app/apperrors"
	"github.com/ervand7/urlshortener/internal/app/database"
	"github.com/ervand7/urlshortener/internal/app/utils"
	q "github.com/ervand7/urlshortener/internal/app/utils/rawqueries"
	"sync"
)

type DBStorage struct {
	DB    database.Database
	Mutex sync.Mutex
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
		return apperrors.NewShortAlreadyExistsError(existsShort)
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
	rows, err := d.DB.Conn.QueryContext(ctx, q.Get, short)
	if err != nil {
		return "", err
	}
	defer d.DB.CloseRows(rows)

	for rows.Next() {
		err = rows.Scan(&origin)
		if err != nil {
			return "", err
		}
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
