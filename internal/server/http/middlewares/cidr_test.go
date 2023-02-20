package middlewares

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrustedNetworkHandler(t *testing.T) {
	type want struct {
		code int
	}
	type header struct {
		name  string
		value string
	}

	tests := []struct {
		name          string
		trustedSubnet string
		header        header
		want          want
	}{
		{
			name:          "access_allowed_x_real_ip",
			trustedSubnet: "192.168.0.1",
			header:        header{name: "X-Real-IP", value: "192.168.0.1"},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          "access_allowed_x_forward_for",
			trustedSubnet: "192.168.0.1",
			header:        header{name: "X-Forwarded-For", value: "192.168.0.1"},
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:          "access_denied",
			trustedSubnet: "192.168.0.1",
			header:        header{name: "X-Real-IP", value: "192.168.0.2"},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name:          "trusted_subnet_empty",
			trustedSubnet: "",
			header:        header{name: "X-Real-IP", value: "192.168.0.1"},
			want: want{
				code: http.StatusForbidden,
			},
		},
		{
			name:          "no_ip_in_headers",
			trustedSubnet: "192.168.0.1",
			want: want{
				code: http.StatusForbidden,
			},
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://test", nil)
			request.Header.Set(tt.header.name, tt.header.value)

			trustedNetwork := &TrustedNetwork{
				TrustedSubnet: net.ParseIP(tt.trustedSubnet),
			}
			handlerToTest := trustedNetwork.Handler(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			)

			result := httptest.NewRecorder()
			handlerToTest.ServeHTTP(result, request)
			response := result.Result()
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
		})
	}
}
