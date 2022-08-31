#! /bin/bash

go test \
  github.com/ervand7/urlshortener/internal/app/config \
  github.com/ervand7/urlshortener/internal/app/controllers/filetable \
  github.com/ervand7/urlshortener/internal/app/controllers/generatedata \
  github.com/ervand7/urlshortener/internal/app/models/url \
  github.com/ervand7/urlshortener/internal/app/server \
  github.com/ervand7/urlshortener/internal/app/views \
  -count 1