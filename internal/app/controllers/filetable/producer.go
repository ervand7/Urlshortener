package filetable

import (
	"bufio"
	"encoding/json"
	"github.com/ervand7/urlshortener/internal/app/config"
	"os"
)

type producer struct {
	file   *os.File
	writer *bufio.Writer
}

func (p *producer) WriteEvent(urlMap map[string]string) error {
	row, err := json.Marshal(urlMap)
	if err != nil {
		return err
	}
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}
	if _, err := p.writer.Write(row); err != nil {
		return err
	}
	return p.writer.Flush()
}

func (p *producer) Close() error {
	return p.file.Close()
}

func NewProducer() (*producer, error) {
	filename := config.GetConfig().FileStoragePath
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}
