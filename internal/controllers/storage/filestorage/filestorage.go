package filestorage

import (
	"context"
	"errors"
	"sync"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
)

type FileStorage struct {
	Mutex sync.Mutex
}

func (f *FileStorage) Set(
	_ context.Context, _ string, short, origin string,
) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	producer, err := newProducer()
	if err != nil {
		return err
	}
	defer func() {
		if err = producer.Close(); err != nil {
			logger.Logger.Warn(err.Error())
		}
	}()

	urlMap := make(map[string]string, 0)
	urlMap[short] = origin
	if err = producer.WriteEvent(urlMap); err != nil {
		return err
	}
	return nil
}

func (f *FileStorage) Get(
	_ context.Context, short string,
) (origin string, err error) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	consumer, err := newConsumer()
	if err != nil {
		return "", err
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			logger.Logger.Warn(err.Error())
		}
	}()

	urlMap, readEventErr := consumer.ReadEvent()
	if readEventErr != nil {
		return "", readEventErr
	}
	origin, exists := urlMap[short]
	if !exists {
		return "", nil
	}
	return origin, nil
}

func (f *FileStorage) SetMany(_ context.Context, _ []models.Entry) error {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return err
}

func (f *FileStorage) GetUserURLs(
	_ context.Context, _ string,
) (result []map[string]string, err error) {
	err = errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return nil, err
}

func (f *FileStorage) DeleteUserURLs(_ []string) {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
}
