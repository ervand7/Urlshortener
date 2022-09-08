#! /bin/bash

go test \
  github.com/ervand7/urlshortener/internal/app/config \
  github.com/ervand7/urlshortener/internal/app/controllers/filesaving \
  github.com/ervand7/urlshortener/internal/app/controllers/generatedata \
  github.com/ervand7/urlshortener/internal/app/server \
  github.com/ervand7/urlshortener/internal/app/views/ \
  -count 1 -v

go test github.com/ervand7/urlshortener/internal/app/models/url/ \
 -run TestFileStorage TestMemoryStorage

export DATABASE_DSN='user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable'
go test github.com/ervand7/urlshortener/internal/app/models/url/ \
  -run TestDBStorage_Set TestDBStorage_Get \
  -count 1 -v