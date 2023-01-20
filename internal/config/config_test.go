package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigDefaultValues(t *testing.T) {
	assert.NotContainsf(t, os.Args, "-a", "flag -a was set")
	assert.Equal(t, os.Getenv("SERVER_ADDRESS"), "")
	expectedServerAddress := ":8080"
	assert.Equal(t, getConfig().ServerAddress, expectedServerAddress)

	assert.NotContainsf(t, os.Args, "-b", "flag -b was set")
	assert.Equal(t, os.Getenv("BASE_URL"), "")
	expectedBaseURL := "http://localhost:8080"
	assert.Equal(t, getConfig().BaseURL, expectedBaseURL)

	assert.NotContainsf(t, os.Args, "-f", "flag -f was set")
	assert.Equal(t, os.Getenv("FILE_STORAGE_PATH"), "")
	expectedFileStoragePath := ""
	assert.Equal(t, getConfig().FileStoragePath, expectedFileStoragePath)
}

func TestGetConfigFromEnv(t *testing.T) {
	serverAddress := ":5000"
	err := os.Setenv("SERVER_ADDRESS", serverAddress)
	assert.NoError(t, err)
	assert.Equal(t, getConfig().ServerAddress, serverAddress)

	baseURL := "http://localhost:777/helloworld"
	err = os.Setenv("BASE_URL", baseURL)
	assert.NoError(t, err)
	assert.Equal(t, getConfig().BaseURL, baseURL)

	fileStoragePath := "/usr/local/go/src/"
	err = os.Setenv("FILE_STORAGE_PATH", fileStoragePath)
	assert.NoError(t, err)
	assert.Equal(t, getConfig().FileStoragePath, fileStoragePath)
}

func TestGetConfigFlagPriority(t *testing.T) {
	serverAddressFlag := ":5000"
	serverAddressEnv := ":7777"
	os.Args = []string{"test", "-a", serverAddressFlag}
	err := os.Setenv("SERVER_ADDRESS", serverAddressEnv)
	assert.NoError(t, err)
	assert.Equal(t, getConfig().ServerAddress, serverAddressFlag)

	BaseURLFlag := "http://hello"
	BaseURLEnv := "http://world"
	os.Args = []string{"test", "-b", BaseURLFlag}
	err = os.Setenv("BASE_URL", BaseURLEnv)
	assert.NoError(t, err)
	assert.Equal(t, getConfig().BaseURL, BaseURLFlag)

	FileStoragePathFlag := "/hello/world/"
	FileStoragePathEnv := "/world/hello/"
	os.Args = []string{"test", "-f", FileStoragePathFlag}
	err = os.Setenv("FILE_STORAGE_PATH", FileStoragePathEnv)
	assert.NoError(t, err)
	assert.Equal(t, getConfig().FileStoragePath, FileStoragePathFlag)
}

func TestGetEnableHTTPS(t *testing.T) {
	value := "hello"
	cacheEnableHTTPS = value
	assert.Equal(t, value, cacheEnableHTTPS)
	assert.Equal(t, value, GetEnableHTTPS())
}

func TestGetServerAddress(t *testing.T) {
	value := "hello"
	cacheServerAddress = value
	assert.Equal(t, value, cacheServerAddress)
	assert.Equal(t, value, GetServerAddress())
}

func TestGetBaseURL(t *testing.T) {
	value := "hello"
	cacheBaseURL = value
	assert.Equal(t, value, cacheBaseURL)
	assert.Equal(t, value, GetBaseURL())
}

func TestGetFileStoragePath(t *testing.T) {
	value := "hello"
	cacheFileStoragePath = value
	assert.Equal(t, value, cacheFileStoragePath)
	assert.Equal(t, value, GetFileStoragePath())
}

func TestGetDatabaseDSN(t *testing.T) {
	value := "hello"
	cacheDatabaseDSN = value
	assert.Equal(t, value, cacheDatabaseDSN)
	assert.Equal(t, value, GetDatabaseDSN())
}
