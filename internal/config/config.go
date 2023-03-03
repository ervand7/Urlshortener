// Package config - app configuration.
package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"

	"github.com/ervand7/urlshortener/internal/logger"
)

// TestDBAddr for tests.
const TestDBAddr = "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable"

var (
	enableHTTPSFlag     *string
	servAddrFlag        *string
	baseURLFlag         *string
	fileStoragePathFlag *string
	databaseDSNFlag     *string
	configFilePath      *string
	trustedSubnetFlag   *string
)

var (
	cacheEnableHTTPS     string
	cacheServerAddress   string
	cacheBaseURL         string
	cacheFileStoragePath string
	cacheDatabaseDSN     string
	cacheTrustedSubnet   string
)

func init() {
	enableHTTPSFlag = flag.String("s", "", "Enable HTTPS")
	servAddrFlag = flag.String("a", "", "Server address")
	baseURLFlag = flag.String("b", "", "Base shorten url")
	fileStoragePathFlag = flag.String("f", "", "File storage path")
	databaseDSNFlag = flag.String("d", "", "Database source name")
	configFilePath = flag.String("c", "", "Config file path")
	trustedSubnetFlag = flag.String("t", "", "trusted subnet addr")
}

type config struct {
	EnableHTTPS     string `env:"ENABLE_HTTPS" json:"enable_https"`
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080" json:"server_address"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	TrustedSubnet   string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`
}

func getConfig() config {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		logger.Logger.Fatal(err.Error())
	}

	flag.Parse()
	if *configFilePath != "" {
		readFromFile(&cfg, *configFilePath)
	}

	if *enableHTTPSFlag == "true" {
		cfg.EnableHTTPS = *enableHTTPSFlag
	}
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
	if *trustedSubnetFlag != "" {
		cfg.TrustedSubnet = *trustedSubnetFlag
	}

	return cfg
}

// GetEnableHTTPS gets EnableHTTPS value by cache.
func GetEnableHTTPS() string {
	if cacheEnableHTTPS != "" {
		return cacheEnableHTTPS
	}
	cacheEnableHTTPS = getConfig().EnableHTTPS
	return cacheEnableHTTPS
}

// GetServerAddress gets serverAddress by cache.
func GetServerAddress() string {
	if cacheServerAddress != "" {
		return cacheServerAddress
	}
	cacheServerAddress = getConfig().ServerAddress
	return cacheServerAddress
}

// GetBaseURL gets baseURL by cache.
func GetBaseURL() string {
	if cacheBaseURL != "" {
		return cacheBaseURL
	}
	cacheBaseURL = getConfig().BaseURL
	return cacheBaseURL
}

// GetFileStoragePath gets fileStoragePath by cache.
func GetFileStoragePath() string {
	if cacheFileStoragePath != "" {
		return cacheFileStoragePath
	}
	cacheFileStoragePath = getConfig().FileStoragePath
	return cacheFileStoragePath
}

// GetDatabaseDSN gets databaseDSN by cache.
func GetDatabaseDSN() string {
	if cacheDatabaseDSN != "" {
		return cacheDatabaseDSN
	}
	cacheDatabaseDSN = getConfig().DatabaseDSN
	return cacheDatabaseDSN
}

// GetTrustedSubnet gets trusted subnet addr.
func GetTrustedSubnet() string {
	if cacheTrustedSubnet != "" {
		return cacheTrustedSubnet
	}
	cacheTrustedSubnet = getConfig().TrustedSubnet
	return cacheTrustedSubnet
}

func readFromFile(cfg *config, path string) {
	configFile, err := os.Open(path)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&cfg)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
