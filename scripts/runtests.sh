#! /bin/bash

go test \
  github.com/ervand7/urlshortener/internal/config \
  github.com/ervand7/urlshortener/internal/controllers/algorithms \
  github.com/ervand7/urlshortener/internal/controllers/storage/filestorage \
  github.com/ervand7/urlshortener/internal/controllers/generatedata \
  github.com/ervand7/urlshortener/internal/server \
  github.com/ervand7/urlshortener/internal/views/ \
  -count 1 -v -p 1

go test github.com/ervand7/urlshortener/internal/controllers/storage/ \
  -count 1 -v -p 1

export DATABASE_DSN='user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable'
go test github.com/ervand7/urlshortener/internal/controllers/storage/ \
  -count 1 -v -p 1
go test github.com/ervand7/urlshortener/internal/views/ \
  -count 1 -v -p 1
