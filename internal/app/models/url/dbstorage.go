package url

import (
	"github.com/ervand7/urlshortener/internal/app/database"
	q "github.com/ervand7/urlshortener/internal/app/utils/rawqueries"
	"sync"
)

type DBStorage struct {
	DB    database.Database
	Mutex sync.Mutex
}

func (d *DBStorage) Set(userID, short, origin string) error {
	_, err := d.DB.Conn.Exec(q.Set, userID, short, origin)
	if err != nil {
		return err
	}

	return nil
}

func (d *DBStorage) Get(short string) (origin string, err error) {
	rows, err := d.DB.Conn.Query(q.Get, short)
	if err != nil {
		return "", err
	}
	defer d.DB.CloseRows(rows)

	err = rows.Scan(&origin)
	if err != nil {
		return "", err
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}

	return origin, nil
}

func (d *DBStorage) GetUserURLs(userID string) (userURLs []map[string]string, err error) {
	rows, err := d.DB.Conn.Query(q.GetUserURLs, userID)
	if err != nil {
		return nil, err
	}
	defer d.DB.CloseRows(rows)

	var (
		short  string
		origin string
	)
	userURLs = make([]map[string]string, 0)
	for rows.Next() {
		err = rows.Scan(&short, &origin)
		if err != nil {
			return nil, err
		}
		pair := map[string]string{
			"short_url":    short,
			"original_url": origin,
		}
		userURLs = append(userURLs, pair)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userURLs, nil
}
