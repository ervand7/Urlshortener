package algorithms

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ervand7/urlshortener/internal/config"
)

func TestGenerateShortURL(t *testing.T) {
	result := GenerateShortURL()
	assert.Contains(t, result, config.GetBaseURL())
	assert.Len(t, result,
		len(config.GetBaseURL())+len("/")+ShortenEndpointLen,
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
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GenerateShortURL()
	}
}

func BenchmarkMakeURLsFromEndpoints(b *testing.B) {
	b.ReportAllocs()
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MakeURLsFromEndpoints(slice)
	}
}
