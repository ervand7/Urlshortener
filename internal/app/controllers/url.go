package controllers

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	BaseUrl           = "http://localhost:8080"
	ShortenUrlLen int = 5
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func ShortenUrl() string {
	result := make([]rune, ShortenUrlLen)
	for i := range result {
		randIndex := rand.Intn(len(letterRunes))
		result[i] = letterRunes[randIndex]
	}
	return BaseUrl + "/" + string(result)
}
