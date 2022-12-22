package db_storage

import (
	"context"
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestNewDBStorage(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer Downgrade()
	db := Database{}
	db.Run()
	storage := NewDBStorage(db)
	assert.Equal(t, storage.db, db)
	assert.Equal(t, storage.resetTimer, true)
	storage.delChan <- "hello"
}

func TestDBStorage_Set(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer Downgrade()
	db := Database{}
	db.Run()
	dbStorage := DBStorage{db: db}

	userID := uuid.New().String()
	short := "http://hello"
	origin := "http://world"

	ctx := context.Background()
	err := dbStorage.Set(ctx, userID, short, origin)
	assert.NoError(t, err)

	type Row struct {
		id     string
		userID string
		short  string
		origin string
		active bool
	}
	rows, err := db.Conn.Query("select * from url")
	assert.NoError(t, err)
	defer db.CloseRows(rows)

	var row Row
	for rows.Next() {
		err = rows.Scan(&row.id, &row.userID, &row.short, &row.origin, &row.active)
		assert.NoError(t, err)
	}

	err = rows.Err()
	assert.NoError(t, err)
	assert.Equal(t, row.userID, userID)
	assert.Equal(t, row.short, short)
	assert.Equal(t, row.origin, origin)
}

func TestDBStorage_SetMany(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer Downgrade()
	db := Database{}
	db.Run()
	dbStorage := DBStorage{db: db}

	userID := uuid.New().String()
	toInsert := []models.Entry{
		{
			UserID: userID,
			Short:  "http://hello1",
			Origin: "http://world1",
		},
		{
			UserID: userID,
			Short:  "http://hello2",
			Origin: "http://world2",
		},
	}
	ctx := context.Background()

	err := dbStorage.SetMany(ctx, toInsert)
	assert.NoError(t, err)

	type Row struct {
		id     string
		userID string
		short  string
		origin string
		active bool
	}
	rows, err := db.Conn.Query("select * from url")
	assert.NoError(t, err)
	defer db.CloseRows(rows)

	entries := make([]Row, 0)
	for rows.Next() {
		row := Row{}
		err = rows.Scan(&row.id, &row.userID, &row.short, &row.origin, &row.active)
		entries = append(entries, row)
		assert.NoError(t, err)
	}

	err = rows.Err()
	assert.NoError(t, err)

	for index := range entries {
		assert.Equal(t, entries[index].userID, toInsert[index].UserID)
		assert.Equal(t, entries[index].short, toInsert[index].Short)
		assert.Equal(t, entries[index].origin, toInsert[index].Origin)
	}
}

func TestDBStorage_Get(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer Downgrade()
	db := Database{}
	db.Run()
	dbStorage := DBStorage{db: db}

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

func TestDBStorage_GetUserURLs(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer Downgrade()
	db := Database{}
	db.Run()
	dbStorage := DBStorage{db: db}

	type UserData struct {
		userID string
		short  string
		origin string
	}
	userID := uuid.New().String()
	toInsert := []UserData{
		{
			userID: userID,
			short:  "hello1",
			origin: "world1",
		},
		{
			userID: userID,
			short:  "hello2",
			origin: "world2",
		},
	}

	ctx := context.Background()
	for _, uData := range toInsert {
		err := dbStorage.Set(ctx, uData.userID, uData.short, uData.origin)
		assert.NoError(t, err)
	}
	result, err := dbStorage.GetUserURLs(ctx, userID)
	assert.NoError(t, err)

	for i := 0; i < len(toInsert); i++ {
		assert.Equal(t, toInsert[i].short, result[i]["short_url"])
		assert.Equal(t, toInsert[i].origin, result[i]["original_url"])
	}
}

func TestDBStorage_DeleteUserURLs(t *testing.T) {
	if os.Getenv("DATABASE_DSN") != config.TestDBAddr {
		return
	}
	defer Downgrade()
	db := Database{}
	db.Run()
	dbStorage := NewDBStorage(db)

	userID := uuid.New().String()
	toInsert := []models.Entry{
		{
			UserID: userID,
			Short:  "http://hello1",
			Origin: "http://world1",
		},
		{
			UserID: userID,
			Short:  "http://hello2",
			Origin: "http://world2",
		},
	}
	ctx := context.Background()
	err := dbStorage.SetMany(ctx, toInsert)
	assert.NoError(t, err)

	dbStorage.DeleteUserURLs([]string{"http://hello1", "http://hello2"})
	time.Sleep(time.Second * Timeout)
	rows, err := db.Conn.QueryContext(ctx, `select "active" from url where "user_id" = $1;`, userID)
	defer db.CloseRows(rows)
	assert.NoError(t, err)

	var active bool
	for rows.Next() {
		err = rows.Scan(&active)
		assert.NoError(t, err)
		assert.Equal(t, false, active)
	}

	err = rows.Err()
	assert.NoError(t, err)
}
