package main

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/ervand7/urlshortener/internal/server/grpc/proto"
)

func getURL(conn *grpc.ClientConn) (*pb.GetURLResponse, error) {
	encoded := hex.EncodeToString([]byte(userID))
	md := metadata.New(map[string]string{"userid": encoded})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := pb.NewShortenerClient(conn)
	respCreate, _ := client.ShortenURL(
		ctx,
		&pb.ShortenURLRequest{
			OriginURL: randomString(60),
		},
	)

	resp, err := client.GetURL(
		ctx,
		&pb.GetURLRequest{
			ShortURL: respCreate.ShortenURL,
		},
	)
	return resp, err
}
