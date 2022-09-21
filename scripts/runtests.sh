#! /bin/bash

go test \
  github.com/ervand7/urlshortener/internal/app/config \
  github.com/ervand7/urlshortener/internal/app/controllers/urlstorage/filestorage \
  github.com/ervand7/urlshortener/internal/app/controllers/generatedata \
  github.com/ervand7/urlshortener/internal/app/server \
  github.com/ervand7/urlshortener/internal/app/views/ \
  -count 1 -v

go test github.com/ervand7/urlshortener/internal/app/controllers/urlstorage/ \
 -run TestFileStorage TestMemoryStorage -count 1 -v

export DATABASE_DSN='user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable'
go test github.com/ervand7/urlshortener/internal/app/controllers/urlstorage/ \
  -run TestDBStorage_Set TestDBStorage_Get \
  -count 1 -v
go test github.com/ervand7/urlshortener/internal/app/views/ \
  -run TestURLShortenBatch \
  -count 1 -v
go test github.com/ervand7/urlshortener/internal/app/views/ \
  -run Test409 \
  -count 1 -v