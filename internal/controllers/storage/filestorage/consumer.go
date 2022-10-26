package filestorage

import (
	"bufio"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/config"
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
			return nil, err
		}
	}
	return urlMap, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func newConsumer() (*consumer, error) {
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
