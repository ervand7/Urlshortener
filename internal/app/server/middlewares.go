package server

import (
	"compress/gzip"
	"github.com/ervand7/urlshortener/internal/app/utils"
	"net/http"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipWrappedBody, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := gzipWrappedBody.Close(); err != nil {
					utils.Logger.Warn(err.Error())
				}
			}()
			r.Body = gzipWrappedBody
		}
		next.ServeHTTP(w, r)
	})
}
