package models

import (
	"context"
)

type Entry struct {
	UserID string
	Short  string
	Origin string
}

type Storage interface {
	Set(ctx context.Context, userID, short, origin string) error
	SetMany(ctx context.Context, dbEntries []Entry) error
	Get(ctx context.Context, short string) (origin string, err error)
	GetUserURLs(ctx context.Context, userID string) (result []map[string]string, err error)
	DeleteUserURLs(shortUrls []string)
}
