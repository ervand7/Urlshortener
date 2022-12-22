package server

import (
	"compress/gzip"
	"net/http"

	"github.com/ervand7/urlshortener/internal/logger"
)

// GzipMiddleware handles gzip body
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipWrappedBody, err := gzip.NewReader(r.Body)
			if err != nil {
				logger.Logger.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				if err = gzipWrappedBody.Close(); err != nil {
					logger.Logger.Warn(err.Error())
				}
			}()
			r.Body = gzipWrappedBody
		}
		next.ServeHTTP(w, r)
	})
}
