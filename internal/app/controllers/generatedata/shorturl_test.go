package generatedata

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestShortenURL(t *testing.T) {
	result := ShortenURL()
	assert.Contains(t, result+"/", config.GetConfig().BaseURL)

	splitResult := strings.Split(result, config.GetConfig().BaseURL+"/")
	endpoint := splitResult[1]
	assert.Equal(t, len(endpoint), ShortenEndpointLen)
}
