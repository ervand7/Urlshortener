package filestorage

import (
	"context"
	"errors"
	"sync"

	"github.com/ervand7/urlshortener/internal/logger"
	"github.com/ervand7/urlshortener/internal/models"
)

// FileStorage implementation of Storage interface for working with files.
type FileStorage struct {
	Mutex sync.Mutex
}

// Set write one entity.
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
		if err = producer.close(); err != nil {
			logger.Logger.Warn(err.Error())
		}
	}()

	urlMap := make(map[string]string, 0)
	urlMap[short] = origin
	if err = producer.writeEvent(urlMap); err != nil {
		return err
	}
	return nil
}

// Get one entity.
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
		if err = consumer.close(); err != nil {
			logger.Logger.Warn(err.Error())
		}
	}()

	urlMap, readEventErr := consumer.readEvent()
	if readEventErr != nil {
		return "", readEventErr
	}
	origin, exists := urlMap[short]
	if !exists {
		return "", nil
	}
	return origin, nil
}

// GetUserURLs gets all user urls from file
func (f *FileStorage) GetUserURLs(
	_ context.Context, _ string,
) (result []map[string]string, err error) {
	err = errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return nil, err
}

// SetMany batch write
func (f *FileStorage) SetMany(_ context.Context, _ []models.Entry) error {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return err
}

// DeleteUserURLs batch delete
func (f *FileStorage) DeleteUserURLs(_ []string) {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
}

// GetNumberOfURLs - get count of shortened urls
func (f *FileStorage) GetNumberOfURLs(ctx context.Context) (int, error) {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return 0, err
}

// GetNumberOfUsers - get count of unique users
func (f *FileStorage) GetNumberOfUsers(ctx context.Context) (int, error) {
	err := errors.New("not implemented")
	logger.Logger.Error(err.Error())
	return 0, err
}
