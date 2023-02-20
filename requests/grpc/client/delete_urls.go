package main

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

func userURLsDelete(conn *grpc.ClientConn) (*pb.UserURLsDeleteResponse, error) {
	encoded := hex.EncodeToString([]byte(userID))
	md := metadata.New(map[string]string{"userid": encoded})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := pb.NewShortenerClient(conn)
	respUserURLs, _ := client.UserURLs(ctx, &pb.UserURLsRequest{})

	var urlsToDelete []string
	for _, val := range respUserURLs.Items {
		urlsToDelete = append(urlsToDelete, val.ShortURL)
	}

	resp, err := client.UserURLsDelete(
		ctx,
		&pb.UserURLsDeleteRequest{
			URLs: urlsToDelete,
		},
	)
	return resp, err
}
