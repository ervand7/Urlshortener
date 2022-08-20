package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"log"
	"os"
)

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (c *consumer) ReadEvent() (map[string]string, error) {
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	urlMap := make(map[string]string, 0)
	for c.scanner.Scan() {
		row := c.scanner.Bytes()
		if err := json.Unmarshal(row, &urlMap); err != nil {
			log.Fatal(err)
		}
	}
	return urlMap, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func NewConsumer() (*consumer, error) {
	filename := config.GetConfig().FileStoragePath
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (f FileTable) Get(key string) (originURL string) {
	consumer, err := NewConsumer()
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	urlMap, readEventErr := consumer.ReadEvent()
	if readEventErr != nil {
		panic(readEventErr)
	}
	originURL, exists := urlMap[key]
	if !exists {
		return ""
	}
	return originURL
}
