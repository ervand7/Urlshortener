package filetable

import (
	"github.com/ervand7/urlshortener/internal/app/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileTable(t *testing.T) {
	fileTable := FileTable{}
	currentDir, currentDirErr := os.Getwd()
	assert.NoError(t, currentDirErr)

	filePath := currentDir + "/file.txt"
	defer os.Remove(filePath)
	os.Args = []string{"test", "-f", filePath}
	assert.Equal(t, config.GetConfig().FileStoragePath, filePath)

	key := "Hello"
	value := "World"
	assert.NoError(t, fileTable.Set(key, value))
	assert.FileExists(t, filePath)

	result, err := fileTable.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, result, value)
}
