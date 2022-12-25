package algorithms

import (
	"math/rand"
	"time"

	"github.com/ervand7/urlshortener/internal/config"
)

const ShortenEndpointLen int = 5

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateShortURL() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]rune, ShortenEndpointLen)
	for i := range result {
		randIndex := rand.Intn(len(letterRunes))
		result[i] = letterRunes[randIndex]
	}
	return config.GetBaseURL() + "/" + string(result)
}

func MakeURLsFromEndpoints(arr []string) {
	baseURL := config.GetBaseURL() + "/"
	for index, val := range arr {
		if len(val) == ShortenEndpointLen {
			arr[index] = baseURL + val
		}
	}
}
