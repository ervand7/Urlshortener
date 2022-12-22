// Package models with interfaces to implement.
package models

import (
	"context"
)

// Entry struct for input param dbEntries in Storage.SetMany.
type Entry struct {
	UserID string
	Short  string
	Origin string
}

// Storage interface for working with storage.
type Storage interface {
	Set(ctx context.Context, userID, short, origin string) error
	SetMany(ctx context.Context, dbEntries []Entry) error
	Get(ctx context.Context, short string) (origin string, err error)
	GetUserURLs(ctx context.Context, userID string) (result []map[string]string, err error)
	DeleteUserURLs(shortUrls []string)
}
