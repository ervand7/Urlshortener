package urlstorage

import (
	"context"
	d "github.com/ervand7/urlshortener/internal/app/database"
	q "github.com/ervand7/urlshortener/internal/app/database/rawqueries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDBStorage_Set(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
		return
	}

	db := d.Database{}
	db.Run()
	dbStorage := DBStorage{DB: db}
	defer func() {
		_, err := db.Conn.Exec(q.DropAll)
		assert.NoError(t, err)
	}()

	userID := uuid.New().String()
	short := "http://hello"
	origin := "http://world"

	ctx := context.Background()
	err := dbStorage.Set(ctx, userID, short, origin)
	assert.NoError(t, err)

	type Entry struct {
		id     string
		userID string
		short  string
		origin string
		active bool
	}
	rows, err := db.Conn.Query("select * from url")
	assert.NoError(t, err)
	defer func() {
		err := rows.Close()
		assert.NoError(t, err)
	}()

	var entry Entry
	for rows.Next() {
		err = rows.Scan(&entry.id, &entry.userID, &entry.short, &entry.origin, &entry.active)
		assert.NoError(t, err)
	}

	err = rows.Err()
	assert.NoError(t, err)
	assert.Equal(t, entry.userID, userID)
	assert.Equal(t, entry.short, short)
	assert.Equal(t, entry.origin, origin)
}

func TestDBStorage_Get(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
		return
	}

	db := d.Database{}
	db.Run()
	dbStorage := DBStorage{DB: db}
	defer func() {
		_, err := db.Conn.Exec(q.DropAll)
		assert.NoError(t, err)
	}()

	userID := uuid.New().String()
	short := "http://hello"
	origin := "http://world"

	ctx := context.Background()
	err := dbStorage.Set(ctx, userID, short, origin)
	assert.NoError(t, err)

	result, err := dbStorage.Get(ctx, short)
	assert.NoError(t, err)
	assert.Equal(t, result, origin)
}
