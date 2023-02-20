package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/logger"
)

// TrustedNetwork for handling trusted network
type TrustedNetwork struct {
	TrustedSubnet net.IP
}

// NewTrustedNetwork constructor of TrustedNetwork
func NewTrustedNetwork() *TrustedNetwork {
	trustedSubnet := net.ParseIP(config.GetTrustedSubnet())
	return &TrustedNetwork{
		TrustedSubnet: trustedSubnet,
	}
}

// Handler checks whether the client's IP address is included in the trusted subnet
func (t *TrustedNetwork) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgProhibited := "Access to the internal network is prohibited"

		if t.TrustedSubnet == nil {
			logger.Logger.Info("empty TrustedSubnet")
			http.Error(w, msgProhibited, http.StatusForbidden)
			return
		}

		ipRaw := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(ipRaw)
		if ip == nil {
			ipRaw = r.Header.Get("X-Forwarded-For")
			ipAddresses := strings.Split(ipRaw, ",")
			ipRaw = ipAddresses[0]
			ip = net.ParseIP(ipRaw)
		}
		if !t.TrustedSubnet.Equal(ip) {
			logger.Logger.Warn(
				fmt.Sprintf("TrustedSubnet - %s, ip - %s", t.TrustedSubnet, ip),
			)
			http.Error(w, msgProhibited, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
