package main

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

func userURLs(conn *grpc.ClientConn) (*pb.UserURLsResponse, error) {
	encoded := hex.EncodeToString([]byte(userID))
	md := metadata.New(map[string]string{"userid": encoded})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := pb.NewShortenerClient(conn)
	resp, err := client.UserURLs(ctx, &pb.UserURLsRequest{})
	return resp, err
}
