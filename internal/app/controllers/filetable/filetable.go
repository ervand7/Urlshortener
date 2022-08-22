package filetable

type FileTable struct{}

func (f FileTable) Get(key string) (originURL string, err error) {
	consumer, err := NewConsumer()
	if err != nil {
		return "", err
	}
	defer consumer.Close()

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
	producer, err := NewProducer()
	if err != nil {
		return err
	}
	defer producer.Close()

	urlMap := make(map[string]string, 0)
	urlMap[key] = value
	if err := producer.WriteEvent(urlMap); err != nil {
		return err
	}
	return nil
}
