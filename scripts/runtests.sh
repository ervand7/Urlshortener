#! /bin/bash

go test \
  github.com/ervand7/urlshortener/internal/config \
  github.com/ervand7/urlshortener/internal/controllers/algorithms \
  github.com/ervand7/urlshortener/internal/controllers/storage/filestorage \
  github.com/ervand7/urlshortener/internal/controllers/storage/memorystorage \
  github.com/ervand7/urlshortener/internal/server/middlewares \
  github.com/ervand7/urlshortener/internal/views/ \
  -count 1 -v -p 10 -bench=. -cpu 8 -benchmem

export DATABASE_DSN='user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable'
go test github.com/ervand7/urlshortener/internal/controllers/storage/dbstorage \
  -count 1 -v -p 1
go test github.com/ervand7/urlshortener/internal/views/ \
  -count 1 -v -p 1
