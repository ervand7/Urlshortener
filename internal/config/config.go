package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/ervand7/urlshortener/internal/logger"
)

const TestDBAddr = "user=ervand password=ervand dbname=urlshortener_test host=localhost port=5432 sslmode=disable"

var (
	servAddrFlag        *string
	baseURLFlag         *string
	fileStoragePathFlag *string
	databaseDSNFlag     *string
)

func init() {
	servAddrFlag = flag.String("a", "", "Server address")
	baseURLFlag = flag.String("b", "", "Base shorten url")
	fileStoragePathFlag = flag.String("f", "", "File storage path")
	databaseDSNFlag = flag.String("d", "", "Database source name")
}

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func GetConfig() Config {
	var cfg Config
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
