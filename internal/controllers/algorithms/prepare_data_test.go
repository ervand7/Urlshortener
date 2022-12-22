package algorithms

import (
	g "github.com/ervand7/urlshortener/internal/controllers/generatedata"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestPrepareShortened(t *testing.T) {
	var (
		slice      []string
		characters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	)
	for i := 0; i < 1000; i++ {
		runes := make([]rune, g.ShortenEndpointLen)
		for i := 0; i < g.ShortenEndpointLen; i++ {
			runes[i] = characters[rand.Intn(len(characters))]
		}
		slice = append(slice, string(runes))
	}

	sliceCopy := make([]string, len(slice))
	copy(sliceCopy, slice)
	PrepareShortened(slice)
	for index, val := range sliceCopy {
		assert.Equal(t, "http://localhost:8080/"+val, slice[index])
	}
}
