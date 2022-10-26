package filestorage

import (
	"context"
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileStorage(t *testing.T) {
	fileStorage := FileStorage{}
	currentDir, currentDirErr := os.Getwd()
	assert.NoError(t, currentDirErr)

	filePath := currentDir + "/file.txt"
	defer func() {
		err := os.Remove(filePath)
		assert.NoError(t, err)
	}()
	os.Args = []string{"test", "-f", filePath}
	assert.Equal(t, config.GetConfig().FileStoragePath, filePath)

	short := "Hello"
	origin := "World"
	assert.NoError(t, fileStorage.Set(context.TODO(), "", short, origin))
	assert.FileExists(t, filePath)

	result, err := fileStorage.Get(context.TODO(), short)
	assert.NoError(t, err)
	assert.Equal(t, result, origin)
}
