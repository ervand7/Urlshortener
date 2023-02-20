package main

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		`:3200`,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resultShorten, err := shorten(conn)
	log.Printf("shortenURL: %v%s\n\n", resultShorten, err)

	resultShortenBatch, err := shortenBatch(conn)
	log.Printf("shortenBatch: %v%s\n\n", resultShortenBatch, err)

	resultGetURL, err := getURL(conn)
	log.Printf("getUrl: %v%s\n\n", resultGetURL, err)

	resultUserURLs, err := userURLs(conn)
	log.Printf("userURLs: %v%s\n\n", resultUserURLs, err)

	resultUserURLsDelete, err := userURLsDelete(conn)
	log.Printf("userURLsDelete: %v%s\n\n", resultUserURLsDelete, err)

	resultGetStats, err := getStats(conn)
	log.Printf("getStats: %#v%s\n\n", resultGetStats, err)
}
