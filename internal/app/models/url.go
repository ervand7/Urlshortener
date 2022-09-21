package models

import (
	"context"
	"github.com/ervand7/urlshortener/internal/app/utils"
)

type Storage interface {
	Set(ctx context.Context, userID, short, origin string) error
	SetMany(ctx context.Context, dbEntries []utils.DBEntry) error
	Get(ctx context.Context, short string) (origin string, err error)
	GetUserURLs(ctx context.Context, userID string) (result []map[string]string, err error)
}
