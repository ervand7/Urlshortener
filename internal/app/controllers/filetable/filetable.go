package filetable

import "github.com/ervand7/urlshortener/internal/app/utils"

type FileTable struct{}

func (f FileTable) Get(key string) (originURL string, err error) {
	consumer, err := newConsumer()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			utils.Logger.Warn(err.Error())
		}
	}()

	urlMap, readEventErr := consumer.ReadEvent()
	if readEventErr != nil {
		return "", readEventErr
	}
	originURL, exists := urlMap[key]
	if !exists {
		return "", nil
	}
	return originURL, nil
}

func (f FileTable) Set(key, value string) error {
	producer, err := newProducer()
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			utils.Logger.Warn(err.Error())
		}
	}()

	urlMap := make(map[string]string, 0)
	urlMap[key] = value
	if err := producer.WriteEvent(urlMap); err != nil {
		return err
	}
	return nil
}
