// Package config - app configuration.
package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/ervand7/urlshortener/internal/logger"
)

// TestDBAddr for tests
const TestDBAddr = "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable"

var (
	servAddrFlag        *string
	baseURLFlag         *string
	fileStoragePathFlag *string
	databaseDSNFlag     *string
)

var (
	cacheServerAddress   string
	cacheBaseURL         string
	cacheFileStoragePath string
	cacheDatabaseDSN     string
)

func init() {
	servAddrFlag = flag.String("a", "", "Server address")
	baseURLFlag = flag.String("b", "", "Base shorten url")
	fileStoragePathFlag = flag.String("f", "", "File storage path")
	databaseDSNFlag = flag.String("d", "", "Database source name")
}

type config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func getConfig() config {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}

	flag.Parse()
	if *servAddrFlag != "" {
		cfg.ServerAddress = *servAddrFlag
	}
	if *baseURLFlag != "" {
		cfg.BaseURL = *baseURLFlag
	}
	if *fileStoragePathFlag != "" {
		cfg.FileStoragePath = *fileStoragePathFlag
	}
	if *databaseDSNFlag != "" {
		cfg.DatabaseDSN = *databaseDSNFlag
	}

	return cfg
}

// GetServerAddress gets serverAddress by cache
func GetServerAddress() string {
	if cacheServerAddress != "" {
		return cacheServerAddress
	}
	cacheServerAddress = getConfig().ServerAddress
	return cacheServerAddress
}

// GetBaseURL gets baseURL by cache
func GetBaseURL() string {
	if cacheBaseURL != "" {
		return cacheBaseURL
	}
	cacheBaseURL = getConfig().BaseURL
	return cacheBaseURL
}

// GetFileStoragePath gets fileStoragePath by cache
func GetFileStoragePath() string {
	if cacheFileStoragePath != "" {
		return cacheFileStoragePath
	}
	cacheFileStoragePath = getConfig().FileStoragePath
	return cacheFileStoragePath
}

// GetDatabaseDSN gets databaseDSN by cache
func GetDatabaseDSN() string {
	if cacheDatabaseDSN != "" {
		return cacheDatabaseDSN
	}
	cacheDatabaseDSN = getConfig().DatabaseDSN
	return cacheDatabaseDSN
}
