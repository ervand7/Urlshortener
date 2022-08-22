package views

import "github.com/ervand7/urlshortener/internal/app/models/url"

type Server struct {
	MemoryStorage *url.MemoryStorage
	FileStorage   *url.FileStorage
}
