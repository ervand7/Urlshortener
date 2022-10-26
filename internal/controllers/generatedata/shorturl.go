package generatedata

import (
	"github.com/ervand7/urlshortener/internal/config"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const ShortenEndpointLen int = 5

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ShortenURL() string {
	result := make([]rune, ShortenEndpointLen)
	for i := range result {
		randIndex := rand.Intn(len(letterRunes))
		result[i] = letterRunes[randIndex]
	}
	return config.GetConfig().BaseURL + "/" + string(result)
}
