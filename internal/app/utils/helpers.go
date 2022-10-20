package utils

type Entry struct {
	UserID string
	Short  string
	Origin string
}

type (
	ReqPair struct {
		CorrelationID string `json:"correlation_id"`
		OriginURL     string `json:"original_url"`
	}
	RespPair struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}
)
