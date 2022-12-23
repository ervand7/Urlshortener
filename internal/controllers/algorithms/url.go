package algorithms

import (
	"math/rand"
	"time"

	"github.com/ervand7/urlshortener/internal/config"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const ShortenEndpointLen int = 5

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateShortURL() string {
	result := make([]rune, ShortenEndpointLen)
	for i := range result {
		randIndex := rand.Intn(len(letterRunes))
		result[i] = letterRunes[randIndex]
	}
	return config.GetConfig().BaseURL + "/" + string(result)
}

func MakeURLsFromEndpoints(arr []string) {
	for index, val := range arr {
		if len(val) == ShortenEndpointLen {
			arr[index] = config.GetConfig().BaseURL + "/" + val
		}
	}
}
