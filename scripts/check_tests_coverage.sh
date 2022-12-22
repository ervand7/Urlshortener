#! /bin/bash

DSN='user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable'

# shellcheck disable=SC1007
DATABASE_DSN=$DSN go test ./... -coverprofile=coverage.out \
  -count 1 -v -p 1 && go tool cover -func coverage.out | grep total:
