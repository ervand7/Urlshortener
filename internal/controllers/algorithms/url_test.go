package algorithms

import (
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	result := GenerateShortURL()
	assert.Contains(t, result, config.GetConfig().BaseURL)
	assert.Len(t, result,
		len(config.GetConfig().BaseURL)+len("/")+ShortenEndpointLen,
	)
}

func TestMakeURLsFromEndpoints(t *testing.T) {
	var (
		slice      []string
		characters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	)
	for i := 0; i < 1000; i++ {
		runes := make([]rune, ShortenEndpointLen)
		for i := 0; i < ShortenEndpointLen; i++ {
			runes[i] = characters[rand.Intn(len(characters))]
		}
		slice = append(slice, string(runes))
	}

	sliceCopy := make([]string, len(slice))
	copy(sliceCopy, slice)
	MakeURLsFromEndpoints(slice)
	for index, val := range sliceCopy {
		assert.Equal(t, "http://localhost:8080/"+val, slice[index])
	}
}

func BenchmarkGenerateShortURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateShortURL()
	}
}

func BenchmarkMakeURLsFromEndpoints(b *testing.B) {
	b.StopTimer()
	var (
		slice      []string
		characters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	)
	for i := 0; i < 1000; i++ {
		runes := make([]rune, ShortenEndpointLen)
		for i := 0; i < ShortenEndpointLen; i++ {
			runes[i] = characters[rand.Intn(len(characters))]
		}
		slice = append(slice, string(runes))
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MakeURLsFromEndpoints(slice)
	}
}
