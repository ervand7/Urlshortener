package filestorage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/ervand7/urlshortener/internal/config"
)

type producer struct {
	file   *os.File
	writer *bufio.Writer
}

func (p *producer) writeEvent(urlMap map[string]string) error {
	row, err := json.Marshal(urlMap)
	if err != nil {
		return err
	}
	if err = p.writer.WriteByte('\n'); err != nil {
		return err
	}
	if _, err = p.writer.Write(row); err != nil {
		return err
	}
	return p.writer.Flush()
}

func (p *producer) close() error {
	return p.file.Close()
}

func newProducer() (*producer, error) {
	filename := config.GetFileStoragePath()
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}
