package filesaving

import "github.com/ervand7/urlshortener/internal/app/utils"

type FileTable struct{}

func (f FileTable) Get(short string) (origin string, err error) {
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
	origin, exists := urlMap[short]
	if !exists {
		return "", nil
	}
	return origin, nil
}

func (f FileTable) Set(short, origin string) error {
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
	urlMap[short] = origin
	if err := producer.WriteEvent(urlMap); err != nil {
		return err
	}
	return nil
}
