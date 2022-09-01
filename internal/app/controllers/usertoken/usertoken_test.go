package usertoken

import (
	"crypto/sha256"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserToken(t *testing.T) {
	userToken := UserToken{Key: sha256.Sum256([]byte("x35k9f"))}
	src := uuid.New().String()

	encoded, err := userToken.Encode(src)
	assert.NoError(t, err)
	assert.NotEqual(t, string(encoded), src)

	decoded, err := userToken.Decode(string(encoded))
	assert.NoError(t, err)
	assert.Equal(t, decoded, src)
}
