package url

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileStorage(t *testing.T) {
	fileStorage := FileStorage{}
	currentDir, currentDirErr := os.Getwd()
	assert.NoError(t, currentDirErr)

	filePath := currentDir + "/file.txt"
	defer os.Remove(filePath)
	os.Args = []string{"test", "-f", filePath}
	assert.Equal(t, config.GetConfig().FileStoragePath, filePath)

	key := "Hello"
	value := "World"
	assert.NoError(t, fileStorage.Set(key, value))
	assert.FileExists(t, filePath)

	result, err := fileStorage.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, result, value)
}
