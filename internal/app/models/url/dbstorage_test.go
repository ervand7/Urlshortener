package url

import (
	"context"
	d "github.com/ervand7/urlshortener/internal/app/database"
	q "github.com/ervand7/urlshortener/internal/app/utils/rawqueries"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDBStorage_Set(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable" {
		return
	}
	d.ManageDB()
	dbStorage := DBStorage{DB: d.DB}
	defer func() {
		_, err := dbStorage.DB.Conn.Exec(q.DropAll)
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
	}
	rows, err := d.DB.Conn.Query("select * from url")
	assert.NoError(t, err)
	defer rows.Close()

	var entry Entry
	for rows.Next() {
		err = rows.Scan(&entry.id, &entry.userID, &entry.short, &entry.origin)
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
	d.ManageDB()
	dbStorage := DBStorage{DB: d.DB}
	defer func() {
		_, err := dbStorage.DB.Conn.Exec(q.DropAll)
		assert.NoError(t, err)
	}()

	userID := uuid.New().String()
	short := "http://hello"
	origin := "http://world"

	ctx := context.Background()
	err := dbStorage.Set(ctx, userID, short, origin)
	assert.NoError(t, err)

	result, err := dbStorage.GetOriginByShort(ctx, short)
	assert.NoError(t, err)
	assert.Equal(t, result, origin)
}
