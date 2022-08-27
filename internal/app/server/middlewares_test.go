package server

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {
	message := "Hello, World!"

	var b bytes.Buffer
	gw := gzip.NewWriter(&b)
	gw.Write([]byte(message))
	gw.Close()

	r := http.Request{
		Body: io.NopCloser(
			bytes.NewReader(
				b.Bytes(),
			),
		),
		Header: map[string][]string{
			"Content-Encoding": {"gzip"},
		},
	}

	wrapped := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decompressed, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		r.Body.Close()
		assert.Equal(t, message, string(decompressed))
	})
	GzipMiddleware(wrapped).ServeHTTP(
		httptest.NewRecorder(), &r,
	)
}
