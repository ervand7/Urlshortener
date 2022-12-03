package algorithms

import (
	"github.com/ervand7/urlshortener/internal/config"
	"github.com/ervand7/urlshortener/internal/controllers/generatedata"
)

func PrepareShortened(arr []string) {
	for index, val := range arr {
		if len(val) == generatedata.ShortenEndpointLen {
			arr[index] = config.GetConfig().BaseURL + "/" + val
		}
	}
}